package worker

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-net/client"
	tcp "github.com/go-pantheon/fabrica-net/tcp/client"
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/xsync"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	"github.com/go-pantheon/roma/gen/app/codec"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/job"
	"github.com/go-pantheon/roma/mercury/internal/job/system"
	"github.com/go-pantheon/roma/mercury/internal/task"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ core.Worker = (*Worker)(nil)

type Worker struct {
	xsync.Stoppable

	logger log.Logger
	log    *log.Helper

	tcpCli   *tcp.Client
	adminCli *core.AdminClients
	dialog   xnet.ClientDialog

	taskChan chan task.Taskable

	userId     int64
	sid        int64
	hsInfo     *HandshakeInfo
	clientUser *climsg.UserProto
	session    xnet.Session

	startTime    time.Time
	workDuration time.Duration

	CompletedChan chan struct{}
	completed     atomic.Bool
	SentMsgCount  atomic.Uint32
	RecvMsgCount  atomic.Uint32
	PushMsgCount  atomic.Uint32
}

func NewWorker(userId int64, logger log.Logger) *Worker {
	w := &Worker{
		Stoppable:     xsync.NewStopper(time.Second * 10),
		logger:        logger,
		userId:        userId,
		sid:           core.ServerId(),
		hsInfo:        &HandshakeInfo{},
		session:       xnet.DefaultSession(),
		taskChan:      make(chan task.Taskable),
		CompletedChan: make(chan struct{}),
	}

	w.log = log.NewHelper(log.With(logger, "module", fmt.Sprintf("worker-%d", userId)))

	return w
}

func (w *Worker) Start(ctx context.Context, jobs []*job.Job) (err error) {
	w.hsInfo, err = newHandshakeInfo(w.userId)
	if err != nil {
		return err
	}

	addrs := core.BootConf().Gate.Addr
	addr := addrs[w.userId%int64(len(addrs))]

	w.tcpCli = tcp.NewClient(w.userId, addr, func(id uint64) (xnet.Pack, error) {
		return w.handshakePack(ctx, w.hsInfo.token, w.hsInfo.cliPub[:])
	}, client.WithAuthFunc(w.Auth))

	if err := w.tcpCli.Start(ctx); err != nil {
		return err
	}

	dialog, ok := w.tcpCli.DefaultDialog()
	if !ok {
		return errors.New("tcp client disconnected")
	}

	w.dialog = dialog

	w.GoAndStop(fmt.Sprintf("worker.%d.work", w.userId), func() error {
		return w.Run(ctx, jobs)
	}, func() error {
		return w.Stop(ctx)
	})

	return nil
}

func (w *Worker) Run(ctx context.Context, jobs []*job.Job) error {
	w.startTime = time.Now()

	w.log.Infof("[worker-%d] connect to gate: %s", w.userId, w.tcpCli.Bind())

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-w.StopTriggered():
			return xsync.ErrStopByTrigger
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		return w.assign(ctx, jobs)
	})
	eg.Go(func() error {
		return w.work(ctx)
	})

	return eg.Wait()
}

func (w *Worker) assign(ctx context.Context, jobs []*job.Job) error {
	dealTicker := time.NewTicker(core.AppConf().WorkMinInterval.AsDuration())
	defer dealTicker.Stop()

	heartbeatTicker := time.NewTicker(core.AppConf().HeartbeatInterval.AsDuration())
	defer heartbeatTicker.Stop()

	for _, j := range jobs {
		for _, t := range j.Tasks {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-heartbeatTicker.C:
				w.taskChan <- system.NewHeartbeatTask()
			case <-dealTicker.C:
				w.taskChan <- t
			}
		}
	}

	close(w.taskChan)

	return nil
}

func (w *Worker) work(ctx context.Context) error {
	return xsync.Run(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case t, ok := <-w.taskChan:
				if !ok {
					w.complete()
					return w.Stop(ctx)
				}

				if err := w.send(t.CSPacket()); err != nil {
					return err
				}

				for {
					done, err := w.receive(ctx, t)
					if err != nil {
						return err
					}

					if done {
						break
					}
				}
			}
		}
	})
}

func (w *Worker) send(pkt *clipkt.Packet) (err error) {
	pkt.Index = w.session.IncreaseCSIndex()

	if pkt.Obj == 0 {
		pkt.Obj = w.UID()
	}

	in, err := proto.Marshal(pkt)
	if err != nil {
		err = errors.Wrapf(err, "Packet marshal failed. uid=%d mod=%d seq=%d", w.UID(), pkt.Mod, pkt.Seq)
		return
	}

	w.SentMsgCount.Add(1)
	task.LogCS(w.log, w.UID(), pkt)

	return w.dialog.Send(in)
}

func (w *Worker) receive(ctx context.Context, t task.Taskable) (done bool, err error) {
	var resp *clipkt.Packet

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-w.StopTriggered():
		return false, xsync.ErrStopByTrigger
	case <-ctx.Done():
		return false, ctx.Err()
	case bytes, ok := <-w.dialog.Receive():
		if !ok {
			return false, errors.Errorf("tcp client disconnected")
		}

		if resp, err = w.decode(bytes); err != nil {
			return false, err
		}
	}

	isPush := codec.IsPushSC(climod.ModuleID(resp.Mod), resp.Seq)
	if isPush {
		w.PushMsgCount.Add(1)
	} else {
		w.RecvMsgCount.Add(1)
	}

	task.LogSC(w.log, w.UID(), resp)

	if isPush {
		return false, nil
	}

	return true, t.Receive(w, resp)
}

func (w *Worker) decode(pack []byte) (p *clipkt.Packet, err error) {
	p = &clipkt.Packet{}

	err = proto.Unmarshal(pack, p)
	if err != nil {
		return nil, errors.Wrapf(err, "Packet unmarshal failed. uid=%d mod=%d seq=%d", w.UID(), p.Mod, p.Seq)
	}

	return p, nil
}

func (w *Worker) Stop(ctx context.Context) error {
	return w.TurnOff(func() (err error) {
		if err = w.tcpCli.Stop(ctx); err != nil {
			return err
		}

		w.workDuration = time.Since(w.startTime)

		w.log.Infof("[worker-%d] closed", w.UID())
		w.log.Infof(w.Output())

		return nil
	})
}

func (w *Worker) SetClientUser(p *climsg.UserProto) {
	w.clientUser = p
	w.log.Infof("[worker-%d] set client user: %s", w.UID(), protojson.Format(p))
}

func (w *Worker) GetClientUser() (*climsg.UserProto, error) {
	if w.clientUser == nil {
		return nil, errors.Errorf("clientUser is nil. uid=%d", w.UID())
	}

	return w.clientUser, nil
}

func (w *Worker) Output() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintf("[worker-%d] ", w.UID()))

	if w.completed.Load() {
		s.WriteString("completed.")
	} else {
		s.WriteString("not completed.")
	}

	s.WriteString(fmt.Sprintf(" spent: %.4fs", w.workDuration.Seconds()))
	s.WriteString(fmt.Sprintf(" sent: %d", w.SentMsgCount.Load()))
	s.WriteString(fmt.Sprintf(" recv: %d", w.RecvMsgCount.Load()))
	s.WriteString(fmt.Sprintf(" push: %d", w.PushMsgCount.Load()))

	return s.String()
}

func (w *Worker) complete() {
	w.completed.Store(true)
}

func (w *Worker) Log() *log.Helper {
	return w.log
}

func (w *Worker) UID() int64 {
	return w.userId
}

func (w *Worker) AdminClient() *core.AdminClients {
	return w.adminCli
}
