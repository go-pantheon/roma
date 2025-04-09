package tcp

import (
	"context"
	"fmt"
	"net"

	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-pantheon/fabrica-kit/ip"
	xnet "github.com/go-pantheon/fabrica-net"
	"github.com/go-pantheon/fabrica-net/conf"
	"github.com/go-pantheon/fabrica-net/internal"
	xcontext "github.com/go-pantheon/fabrica-net/xcontext"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/pkg/errors"
)

var _ transport.Server = (*Server)(nil)

type Option func(o *Server)

type WrapFunc func(ctx context.Context, color string, uid int64) error

func Bind(bind string) Option {
	return func(s *Server) {
		s.conf.Server.Bind = bind
	}
}

func Referer(referer string) Option {
	return func(s *Server) {
		s.referer = referer
	}
}

func Logger(logger log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func ReadFilter(m middleware.Middleware) Option {
	return func(s *Server) {
		if s.readFilter == nil {
			s.readFilter = m
			return
		}
		s.readFilter = middleware.Chain(s.readFilter, m)
	}
}

func WriteFilter(m middleware.Middleware) Option {
	return func(s *Server) {
		if s.writeFilter == nil {
			s.writeFilter = m
			return
		}
		s.writeFilter = middleware.Chain(s.writeFilter, m)
	}
}

func AfterConnectFunc(f WrapFunc) Option {
	return func(s *Server) {
		s.afterConnectFunc = f
	}
}

func AfterDisconnectFunc(f WrapFunc) Option {
	return func(s *Server) {
		s.afterDisconnectFunc = f
	}
}

type Server struct {
	xsync.Stoppable

	conf    *conf.Config
	logger  log.Logger
	referer string

	workerSize int
	listener   net.Listener
	buckets    *internal.Buckets

	handler     xnet.Service
	readFilter  middleware.Middleware
	writeFilter middleware.Middleware

	afterConnectFunc    WrapFunc
	afterDisconnectFunc WrapFunc
}

func NewServer(handler xnet.Service, opts ...Option) (*Server, error) {
	conf.Init()

	s := &Server{
		Stoppable: xsync.NewStopper(conf.Conf.Server.StopTimeout),
		conf:      conf.Conf,
		logger:    log.DefaultLogger,
		readFilter: middleware.Chain(
			recovery.Recovery(),
		),
		writeFilter: middleware.Chain(
			recovery.Recovery(),
		),
		handler: handler,
	}

	for _, o := range opts {
		o(s)
	}

	s.buckets = internal.NewBuckets(s.conf.Bucket)
	s.workerSize = s.conf.Server.WorkerSize

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	var (
		listener *net.TCPListener
		addr     *net.TCPAddr
		err      error
	)

	bind := s.conf.Server.Bind
	if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
		err = errors.Wrapf(err, "resolve bind failed. bind=%s", bind)
		return err
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		err = errors.Wrapf(err, "listen failed. addr=%s", addr.String())
		return err
	}

	xcontext.SetDeadlineWithContext(ctx, listener, "TcpListener")

	s.listener = listener
	idGen := &atomic.Uint64{}
	for i := 0; i < s.workerSize; i++ {
		workerID := i
		xsync.GoSafe(fmt.Sprintf("tcp.Server.acceptLoop.%d", workerID), func() error {
			return s.acceptLoop(ctx, idGen)
		})
	}

	log.Infof("[tcp.Server] listening on %s", addr.String())
	return nil
}

func (s *Server) Stop(ctx context.Context) (err error) {
	s.stop()
	return
}

func (s *Server) stop() {
	s.DoStop(func() {
		s.buckets.Walk(func(w *internal.Worker) (continued bool) {
			w.TriggerStop()
			return true
		})
		s.buckets.Walk(func(w *internal.Worker) (continued bool) {
			w.WaitStopped()
			return true
		})
	})
	log.Info("[tcp.Server] TCP server is closed")
}

func (s *Server) acceptLoop(ctx context.Context, idGen *atomic.Uint64) error {
	for {
		select {
		case <-s.Stopping():
			s.WaitStopped()
			return ctx.Err()
		case <-ctx.Done():
			s.WaitStopped()
			return ctx.Err()
		default:
			if err := s.accept(ctx, idGen); err != nil {
				log.Errorf("[tcp.Server] %+v", err)
			}
		}
	}
}

func (s *Server) accept(ctx context.Context, idGen *atomic.Uint64) error {
	conn, err := s.listener.(*net.TCPListener).AcceptTCP()
	if err != nil {
		return errors.Wrapf(err, "accept failed")
	}

	conn0 := conn
	wid := idGen.Add(1)
	xsync.GoSafe(fmt.Sprintf("tcp.Server.serve.%d", wid), func() error {
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		if err := s.serve(ctx, conn0, wid); err != nil {
			return errors.WithMessagef(err, "serve failed wid=%d remote=%s local=%s",
				wid, xcontext.RemoteAddr(conn0), xcontext.LocalAddr(conn0))
		}
		return nil
	})
	return nil
}

func (s *Server) serve(ctx context.Context, conn *net.TCPConn, wid uint64) error {
	if err := conn.SetKeepAlive(s.conf.Server.KeepAlive); err != nil {
		return errors.Wrapf(err, "SetKeepAlive failed v=%v	", s.conf.Server.KeepAlive)
	}
	if err := conn.SetReadBuffer(s.conf.Server.ReadBufSize); err != nil {
		return errors.Wrapf(err, "SetReadBuffer failed v=%d", s.conf.Server.ReadBufSize)
	}
	if err := conn.SetWriteBuffer(s.conf.Server.WriteBufSize); err != nil {
		return errors.Wrapf(err, "SetWriteBuffer failed v=%d", s.conf.Server.WriteBufSize)
	}

	return s.work(ctx, conn, wid)
}

func (s *Server) work(ctx context.Context, conn *net.TCPConn, wid uint64) (err error) {
	w := internal.NewWorker(wid, conn, s.logger, s.conf.Worker, s.referer, s.readFilter, s.writeFilter, s.handler)

	defer func() {
		if err != nil {
			err = errors.WithMessagef(err, "uid=%d color=%s state=%d", w.UID(), w.Color(), w.Status())
		}

		w.Stop(ctx)
		if s.afterDisconnectFunc != nil {
			if err = s.afterDisconnectFunc(ctx, w.Color(), w.UID()); err != nil {
				log.Errorf("[tcp.Server] afterDisconnectFunc failed. wid=%d remote=%s local=%s uid=%d color=%s state=%d %+v",
					w.WID(), xcontext.RemoteAddr(w.Conn()), xcontext.LocalAddr(w.Conn()), w.UID(), w.Color(), w.Status(), err)
			}
		}
	}()

	if err = w.Start(ctx); err != nil {
		return err
	}
	if err = s.putBucket(w); err != nil {
		return err
	}
	defer s.buckets.Del(w)

	if s.afterConnectFunc != nil {
		if err = s.afterConnectFunc(ctx, w.Color(), w.UID()); err != nil {
			return err
		}
	}

	return w.Run(ctx)
}

func (s *Server) putBucket(w *internal.Worker) error {
	if ow := s.buckets.Put(w); ow != nil {
		log.Errorf("[tcp.Server] worker is replaced, close old worker. wid=%d remote=%s local=%s uid=%d color=%s "+
			"old-remote=%s old-local=%s old-color=%s",
			w.WID(), xcontext.RemoteAddr(w.Conn()), xcontext.LocalAddr(w.Conn()), w.UID(), w.Color(),
			xcontext.RemoteAddr(ow.Conn()), xcontext.LocalAddr(ow.Conn()), ow.Color())
		ow.TriggerStop()
	}
	return nil
}

func (s *Server) Disconnect(ctx context.Context, key uint64) error {
	w := s.buckets.Worker(key)
	if w == nil {
		return errors.New("worker not found")
	}

	w.TriggerStop()
	w.WaitStopped()
	return nil
}

func (s *Server) WIDList() []uint64 {
	ids := make([]uint64, 0, 1024)
	s.buckets.Walk(func(w *internal.Worker) bool {
		ids = append(ids, w.WID())
		return true
	})
	return ids
}

func (s *Server) Push(ctx context.Context, uid int64, pack []byte) error {
	if len(pack) <= 0 {
		return errors.New("push msg len <= 0")
	}

	w := s.buckets.GetByUID(uid)
	if w == nil {
		return errors.Errorf("worker not found. uid=%d", uid)
	}
	return w.Push(ctx, pack)
}

func (s *Server) PushGroup(ctx context.Context, uids []int64, pack []byte) (err error) {
	if len(pack) <= 0 {
		return errors.New("push group msg len <= 0")
	}

	workers := s.buckets.GetByUIDs(uids)
	for _, w := range workers {
		if err0 := w.Push(ctx, pack); err0 != nil {
			err = errors.WithMessagef(err0, " uid=%d", w.UID())
		}
	}
	return
}

func (s *Server) Broadcast(ctx context.Context, pack []byte) (err error) {
	if len(pack) <= 0 {
		return errors.New("broadcast msg len <= 0")
	}

	s.buckets.Walk(func(w *internal.Worker) bool {
		if err0 := w.Push(ctx, pack); err0 != nil {
			err = errors.WithMessagef(err0, " uid=%d", w.UID())
		}
		return true
	})
	return
}

func (s *Server) Endpoint() (string, error) {
	addr, err := ip.Extract(s.conf.Server.Bind, s.listener)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("tcp://%s", addr), nil
}
