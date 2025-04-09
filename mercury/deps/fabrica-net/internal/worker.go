package internal

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-pantheon/fabrica-kit/tunnel"
	xnet "github.com/go-pantheon/fabrica-net"
	"github.com/go-pantheon/fabrica-net/conf"
	"github.com/go-pantheon/fabrica-net/internal/bufreader"
	"github.com/go-pantheon/fabrica-net/xcontext"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var _ tunnel.Holder = (*Worker)(nil)
var _ xsync.Stoppable = (*Worker)(nil)

type Worker struct {
	*tunnelHolder
	xsync.Stoppable
	xsync.CountdownStopper

	conf             *conf.Worker
	reader           *bufreader.Reader
	service          xnet.Service
	createTunnelFunc CreateTunnelFunc
	referer          string

	readFilter  middleware.Middleware
	writeFilter middleware.Middleware

	id      uint64
	conn    *net.TCPConn
	started atomic.Bool
	session xnet.Session

	replyChanStarted   atomic.Bool
	replyChanCompleted chan struct{}
	replyChan          chan []byte
}

func NewWorker(wid uint64, conn *net.TCPConn, logger log.Logger, conf *conf.Worker, referer string,
	readFilter, writeFilter middleware.Middleware, handler xnet.Service) *Worker {
	w := &Worker{
		tunnelHolder:       newTunnelHolder(),
		Stoppable:          xsync.NewStopper(conf.StopTimeout),
		CountdownStopper:   xsync.NewCountdownStopper(),
		conf:               conf,
		service:            handler,
		referer:            referer,
		readFilter:         readFilter,
		writeFilter:        writeFilter,
		id:                 wid,
		conn:               conn,
		session:            xnet.DefaultSession(),
		replyChanCompleted: make(chan struct{}),
	}

	w.createTunnelFunc = func(ctx context.Context, tp int32, oid int64) (tunnel.Tunnel, error) {
		t, err := w.service.CreateTunnel(ctx, w.session, tp, oid, w)
		if err != nil {
			return nil, err
		}
		return t, nil
	}

	w.replyChan = make(chan []byte, conf.ReplyChanSize)
	w.reader = bufreader.NewReader(conn, conf.ReaderBufSize)
	return w
}

func (w *Worker) Start(ctx context.Context) (err error) {
	if err = w.Conn().SetDeadline(time.Now().Add(w.conf.HandshakeTimeout)); err != nil {
		return errors.Wrap(err, "set conn deadline before handshake failed")
	}

	if err = w.handshake(ctx); err != nil {
		return err
	}
	if err = w.service.OnConnected(ctx, w.session); err != nil {
		return err
	}

	if err = w.Conn().SetDeadline(time.Now().Add(w.conf.RequestIdleTimeout)); err != nil {
		return errors.Wrap(err, "set conn deadline after handshake failed")
	}

	w.started.Store(true)
	return
}

func (w *Worker) Run(ctx context.Context) error {
	ctx = xcontext.SetUID(ctx, w.UID())
	ctx = xcontext.SetSID(ctx, w.SID())
	ctx = xcontext.SetColor(ctx, w.Color())
	ctx = xcontext.SetStatus(ctx, w.Status())
	ctx = xcontext.SetGateReferer(ctx, w.referer, w.WID())
	ctx = xcontext.SetClientIP(ctx, w.session.ClientIP())

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-w.StopTriggered():
			return xsync.GroupStopping
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		// sc loop
		err := xsync.RunSafe(func() error {
			return w.writePackLoop(ctx)
		})
		return err
	})
	eg.Go(func() error {
		// cs loop
		err := xsync.RunSafe(func() error {
			return w.readPackLoop(ctx)
		})
		return err
	})
	eg.Go(func() error {
		err := w.tickStopSign(ctx)
		return err
	})

	if err := eg.Wait(); err != nil {
		return errors.WithMessagef(err, "uid=%d color=%s", w.UID(), w.Color())
	}
	return nil
}

func (w *Worker) Stop(ctx context.Context) {
	w.DoStop(func() {
		if w.IsStarted() {
			ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			if err := w.service.OnDisconnect(ctx, w.session); err != nil {
				log.Errorf("[xnet.Worker] onDisconnect failed. wid=%d uid=%d color=%s %+v", w.WID(), w.UID(), w.Color(), err)
			}
		}

		close(w.replyChan)
		if w.replyChanStarted.Load() {
			<-w.replyChanCompleted
		}

		if err := w.reader.Close(); err != nil {
			log.Errorf("[xnet.Worker] reader close failed. wid=%d uid=%d color=%s %+v", w.WID(), w.UID(), w.Color(), err)
		}

		w.tunnelHolder.stop()

		if err0 := w.conn.Close(); err0 != nil {
			log.Errorf("[xnet.Worker] conn close failed. wid=%d uid=%d color=%s %+v", w.WID(), w.UID(), w.Color(), err0)
			xcontext.SetDeadlineWithContext(ctx, w.conn, fmt.Sprintf("wid=%d", w.WID()))
		}
	})
}

// handshake must only be used in auth
func (w *Worker) handshake(ctx context.Context) error {
	var (
		ss  xnet.Session
		in  []byte
		out []byte
		err error
	)

	if in, err = w.read(); err != nil {
		return err
	}
	if out, ss, err = w.service.Auth(ctx, in); err != nil {
		return err
	}
	if err = w.write(out); err != nil {
		return err
	}

	ss.SetClientIP(xcontext.RemoteAddr(w.conn))
	w.session = ss
	return nil
}

func (w *Worker) Tunnel(ctx context.Context, mod int32, oid int64) (t tunnel.Tunnel, err error) {
	tp, err := w.service.TunnelType(mod)
	if err != nil {
		return
	}
	if t = w.tunnel(tp, oid); t != nil {
		return
	}
	return w.createTunnel(ctx, tp, oid, w.createTunnelFunc)
}

func (w *Worker) Push(ctx context.Context, out []byte) error {
	if w.IsStopping() {
		return errors.New("worker is stopping")
	}
	if len(out) <= 0 {
		return errors.New("push msg len <= 0")
	}

	w.replyChan <- out
	return nil
}

func (w *Worker) tickStopSign(ctx context.Context) (err error) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if t := w.CountdownStopper.ExpiryTime(); !t.IsZero() && time.Now().After(t) {
				w.TriggerStop()
				return errors.Wrapf(xsync.ErrCountdownTimerExpired, "wid=%d", w.WID())
			}
			// TODO: check black list
		}
	}
}

func (w *Worker) writePackLoop(ctx context.Context) (err error) {
	defer close(w.replyChanCompleted)

	w.replyChanStarted.Store(true)
	for pack := range w.replyChan {
		if err = w.writePack(ctx, pack); err != nil {
			return err
		}
	}
	return nil
}

func (w *Worker) readPackLoop(ctx context.Context) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.Stopping():
			return nil
		default:
			if err = w.readPack(ctx); err != nil {
				return err
			}
		}
	}
}

func writeNext(ctx context.Context, pk interface{}) (interface{}, error) {
	return pk, nil
}

func (w *Worker) writePack(ctx context.Context, pack []byte) (err error) {
	next := writeNext
	if w.writeFilter != nil {
		next = w.writeFilter(next)
	}

	var out interface{}
	if out, err = next(ctx, pack); err != nil {
		return
	}
	return w.write(out.([]byte))
}

func (w *Worker) write(pack []byte) (err error) {
	pack, err = encrypt(w.session, pack)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	err = binary.Write(&buf, binary.BigEndian, uint32(len(pack)))
	if err != nil {
		err = errors.Wrap(err, "write packet length failed")
		return
	}

	if _, err = buf.Write(pack); err != nil {
		err = errors.Wrap(err, "write packet body failed")
		return
	}

	if _, err = w.conn.Write(buf.Bytes()); err != nil {
		err = errors.Wrap(err, "send packet failed")
		return
	}
	return
}

func (w *Worker) readPack(ctx context.Context) (err error) {
	var in []byte
	if in, err = w.read(); err != nil {
		return
	}

	next := func(ctx context.Context, req interface{}) (interface{}, error) {
		return w.handle(ctx, req)
	}
	if w.readFilter != nil {
		next = w.readFilter(next)
	}

	if _, err = next(ctx, in); err != nil {
		return
	}

	_ = w.conn.SetDeadline(time.Now().Add(w.conf.RequestIdleTimeout))
	return
}

func (w *Worker) read() (buf []byte, err error) {
	var lenBytes []byte
	if lenBytes, err = w.reader.ReadFull(xnet.PackLenSize); err != nil {
		err = errors.Wrap(err, "read packet length failed")
		return
	}

	var packLen int32
	if err = binary.Read(bytes.NewReader(lenBytes), binary.BigEndian, &packLen); err != nil {
		return
	}
	if packLen <= 0 {
		err = errors.New("packet len must greater than 0")
		return
	}
	if packLen > xnet.MaxBodySize {
		err = errors.Errorf("packet len=%d must less than %d", packLen, xnet.MaxBodySize)
		return
	}

	if buf, err = w.reader.ReadFull(int(packLen)); err != nil {
		err = errors.Wrapf(err, "read packet body failed. len=%d", packLen)
		return
	}

	if buf, err = decrypt(w.session, buf); err != nil {
		return
	}
	return
}

func (w *Worker) handle(ctx context.Context, req interface{}) (interface{}, error) {
	err := w.service.Handle(ctx, w.session, w, req.([]byte))
	return nil, err
}

func (w *Worker) IsStarted() bool {
	return w.started.Load()
}

// SetStopCountDownTime Pass in the current time to set the worker shutdown countdown when the main tunnel is disconnected
func (w *Worker) SetStopCountDownTime(now time.Time) {
	w.CountdownStopper.SetExpiryTime(now.Add(w.conf.WaitMainTunnelTimeout))
}

func (w *Worker) Conn() *net.TCPConn {
	return w.conn
}

func (w *Worker) Session() xnet.Session {
	return w.session
}

func (w *Worker) WID() uint64 {
	return w.id
}

func (w *Worker) UID() int64 {
	if w.session == nil {
		return 0
	}
	return w.session.UID()
}

func (w *Worker) SID() int64 {
	if w.session == nil {
		return 0
	}
	return w.session.SID()
}

func (w *Worker) Color() string {
	if w.session == nil {
		return ""
	}
	return w.session.Color()
}

func (w *Worker) Status() int64 {
	if w.session == nil {
		return 0
	}
	return w.session.Status()
}

func (w *Worker) Endpoint() string {
	if w.conn == nil {
		return ""
	}
	return w.conn.RemoteAddr().String()
}
