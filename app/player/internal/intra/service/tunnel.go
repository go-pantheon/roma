package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/player/intra/v1"
	"github.com/go-pantheon/roma/gen/app/player/handler"
	"github.com/go-pantheon/roma/gen/app/player/service"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/go-pantheon/roma/pkg/universe/middleware/logging"
	"golang.org/x/sync/errgroup"
)

type TunnelService struct {
	intrav1.UnimplementedTunnelServiceServer

	log *log.Helper
	mgr *core.Manager
	svc *service.PlayerServices
}

func NewTunnelService(logger log.Logger, mgr *core.Manager, svc *service.PlayerServices) intrav1.TunnelServiceServer {
	return &TunnelService{
		log: log.NewHelper(log.With(logger, "module", "player/intra/tunnel")),
		mgr: mgr,
		svc: svc,
	}
}

func (s *TunnelService) Tunnel(stream intrav1.TunnelService_TunnelServer) (err error) {
	defer func() {
		if err != nil {
			s.log.Errorf("tunnel error: %+v", err)
		}
	}()

	ctx := stream.Context()

	if !life.IsGateContext(ctx) {
		return errors.Errorf("must be called by Gateway. status=%d", xcontext.Status(ctx))
	}

	oid, err := xcontext.OID(ctx)
	if err != nil {
		return err
	}

	sendFunc := func(p xnet.TunnelMessage) error {
		return core.SendFunc(stream, p)
	}

	w, err := s.mgr.Worker(ctx, oid, core.NewResponser(sendFunc))
	if err != nil {
		return err
	}

	defer func() {
		w.Stop(ctx)
	}()

	return s.run(ctx, w, stream)
}

func (s *TunnelService) run(ctx context.Context, w life.Workable, stream intrav1.TunnelService_TunnelServer) (err error) {
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
		return xsync.Run(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
					if err = s.recv(w, stream); err != nil {
						return err
					}
				}
			}
		})
	})

	return eg.Wait()
}

func (s *TunnelService) recv(w life.Workable, stream intrav1.TunnelService_TunnelServer) (err error) {
	in, err := stream.Recv()
	if err != nil {
		return errors.Wrapf(err, "intrav1.TunnelRequest recv failed")
	}

	if in == nil {
		return nil
	}

	return w.EmitEventFunc(func(wctx life.Context) error {
		return s.handle(wctx.(core.Context), in)
	})
}

func (s *TunnelService) handle(wctx core.Context, in *intrav1.TunnelRequest) error {
	logging.Req(wctx, s.log, in, logging.DefaultFilter)

	st := time.Now()

	out, err := handler.PlayerHandle(wctx, s.svc, in)
	if err != nil {
		return err
	}

	if out == nil {
		out, err = s.handleError(in)
		if err != nil {
			return err
		}
	}

	logging.Resp(wctx, s.log, wctx.UID(), out, time.Since(st), logging.DefaultFilter)

	wctx.Reply(out)

	return nil
}

func (s *TunnelService) handleError(in *intrav1.TunnelRequest) (xnet.TunnelMessage, error) {
	return handler.TakeProtoPlayerTunnelResponse(
		in.Index,
		int32(climod.ModuleID_System),
		int32(cliseq.SystemSeq_ServerUnexpectedErr),
		in.Obj,
		&climsg.SCServerUnexpectedErr{
			Mod: in.Mod,
			Seq: in.Seq,
		},
	)
}
