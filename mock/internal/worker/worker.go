package worker

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	clipkt "github.com/vulcan-frame/vulcan-game/gen/api/client/packet"
	"github.com/vulcan-frame/vulcan-game/gen/app/codec"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base/security"
	"github.com/vulcan-frame/vulcan-game/mock/internal/job/system"
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
	tcp "github.com/vulcan-frame/vulcan-net/tcp/client"
	"github.com/vulcan-frame/vulcan-util/xsync"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
)

var _ base.UserManager = (*Worker)(nil)

type Worker struct {
	xsync.Stoppable

	crypto *security.Crypto

	logger log.Logger
	log    *log.Helper

	userId int64
	Token  string

	tcpCli   *tcp.Client
	adminCli *base.AdminClients

	taskChan     chan task.Taskable
	redirectChan chan *clipkt.Packet

	clientUser *climsg.UserProto

	Completed    atomic.Bool
	WorkingTime  time.Duration
	SentMsgCount atomic.Uint32
	RecvMsgCount atomic.Uint32
	PushMsgCount atomic.Uint32
}

func (w *Worker) Log() *log.Helper {
	return w.log
}

func (w *Worker) UID() int64 {
	return w.userId
}

func (w *Worker) AdminClient() *base.AdminClients {
	return w.adminCli
}

func NewWorker(userId int64, logger log.Logger) *Worker {
	w := &Worker{
		Stoppable:    xsync.NewStopper(time.Second * 10),
		crypto:       &security.Crypto{},
		logger:       logger,
		userId:       userId,
		taskChan:     make(chan task.Taskable),
		redirectChan: make(chan *clipkt.Packet, 1024),
	}

	w.log = log.NewHelper(log.With(logger, "module", fmt.Sprintf("worker-%d", userId)))

	return w
}

func (w *Worker) Start(ctx *base.Context) (err error) {
	addrs := base.Dep.Conf.Boot.Gate.Addr
	addr := addrs[w.userId%int64(len(addrs))]

	opts := []tcp.Option{
		tcp.Bind(addr),
	}

	w.tcpCli = tcp.NewClient(w.userId, opts...)

	w.Token, err = genToken(w.userId)
	if err != nil {
		return err
	}

	if err = w.tcpCli.Start(ctx); err != nil {
		return err
	}

	w.log.Infof("worker-%d connect to gate: %s", w.userId, addr)

	xsync.GoSafe("worker.start", func() error {
		return w.Work(ctx)
	}, DontLog)
	return nil
}

func (w *Worker) SetClientUser(p *climsg.UserProto) {
	w.clientUser = p
}

func (w *Worker) GetClientUser() (*climsg.UserProto, error) {
	if w.clientUser == nil {
		return nil, errors.Errorf("clientUser is nil. uid=%d", w.UID())
	}

	return w.clientUser, nil
}

func (w *Worker) Stop() {
	w.stop()
	w.log.Infof(w.Output())
}

func (w *Worker) stop() {
	w.DoStop(func() {
		w.tcpCli.TriggerStop()
		w.tcpCli.WaitStopped()
		close(w.redirectChan)

		w.log.Debugf("worker-%d closed", w.UID())
	})
}

func (w *Worker) DistributeTask(t task.Taskable) {
	w.taskChan <- t
}

func (w *Worker) Work(bctx *base.Context) error {
	heartbeatTicker := time.NewTicker(base.App().HeartbeatInterval.AsDuration())
	heartbeat := system.NewHeartbeatTask()

	defer w.stop()

	eg, ctx := errgroup.WithContext(bctx)
	eg.Go(func() error {
		select {
		case <-w.StopTriggered():
			return xsync.GroupStopping
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		return xsync.RunSafe(func() error {
			for {
				select {
				case t := <-w.taskChan:
					if err := w.work(bctx, t); err != nil {
						return err
					}
				case <-heartbeatTicker.C:
					if err := w.work(bctx, heartbeat); err != nil {
						return err
					}
				case redirect := <-w.redirectChan:
					if err := w.send(bctx, redirect); err != nil {
						return err
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	})

	return eg.Wait()
}

func (w *Worker) work(ctx *base.Context, t task.Taskable) error {
	if err := w.send(ctx, t.CSPacket()); err != nil {
		return err
	}

	for {
		if w.IsStopping() {
			return nil
		}
		done, err := w.receive(ctx, t)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

func (w *Worker) receive(ctx *base.Context, t task.Taskable) (done bool, err error) {
	var (
		resp     *clipkt.Packet
		redirect *clipkt.Packet
	)

	timeout := time.NewTimer(5 * time.Second)
	defer timeout.Stop()

	select {
	case <-w.StopTriggered():
		err = xsync.GroupStopping
		return
	case <-timeout.C:
		err = errors.Errorf("worker receive response timeout")
		return
	case bytes, ok := <-w.tcpCli.Receive():
		if !ok {
			err = errors.Errorf("tcp client disconnected")
			return
		}
		if resp, err = w.decode(bytes); err != nil {
			return
		}
	}

	if codec.IsPushSC(climod.ModuleID(resp.Mod), resp.Seq) {
		w.PushMsgCount.Add(1)
	} else {
		w.RecvMsgCount.Add(1)
	}

	task.LogSC(w.log, resp)

	redirect, done, err = t.Receive(ctx, resp)
	if redirect != nil {
		w.redirectChan <- redirect
	}
	return
}

func (w *Worker) send(ctx *base.Context, pkt *clipkt.Packet) (err error) {
	pkt.Index = ctx.SendIndex
	ctx.SendIndex++
	if pkt.Obj == 0 {
		pkt.Obj = w.UID()
	}

	in, err := proto.Marshal(pkt)
	if err != nil {
		err = errors.Wrapf(err, "Packet marshal failed. uid=%d mod=%d seq=%d", w.UID(), pkt.Mod, pkt.Seq)
		return
	}

	if !base.Unencrypted() {
		if in, err = w.crypto.EncryptProto(in); err != nil {
			err = errors.Wrapf(err, "Packet encrypt failed. uid=%d mod=%d seq=%d", w.UID(), pkt.Mod, pkt.Seq)
			return
		}
	}

	w.SentMsgCount.Add(1)
	task.LogCS(w.log, pkt)

	return w.tcpCli.Send(in)
}

func (w *Worker) decode(pack []byte) (body *clipkt.Packet, err error) {
	if !base.Unencrypted() {
		if pack, err = w.crypto.DecryptProto(pack); err != nil {
			return
		}
	}

	p := &clipkt.Packet{}
	err = proto.Unmarshal(pack, p)
	if err != nil {
		err = errors.Wrapf(err, "Packet unmarshal failed. uid=%d mod=%d seq=%d", w.UID(), p.Mod, p.Seq)
		return
	}
	return
}

func (w *Worker) Output() string {
	s := &strings.Builder{}
	s.WriteString(fmt.Sprintf("worker-%d ", w.UID()))
	if w.Completed.Load() {
		s.WriteString("completed")
	} else {
		s.WriteString("not completed")
	}
	s.WriteString(fmt.Sprintf(" time: %.4fs", w.WorkingTime.Seconds()))
	s.WriteString(fmt.Sprintf(" sent: %d", w.SentMsgCount.Load()))
	s.WriteString(fmt.Sprintf(" recv: %d", w.RecvMsgCount.Load()))
	s.WriteString(fmt.Sprintf(" push: %d", w.PushMsgCount.Load()))

	return s.String()
}
