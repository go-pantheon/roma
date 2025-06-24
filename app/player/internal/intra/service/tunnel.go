package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xcontext"
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
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
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

func (s *TunnelService) Tunnel(stream intrav1.TunnelService_TunnelServer) error {
	ctx := stream.Context()

	if !life.IsGateContext(ctx) {
		return errors.Errorf("must be called by Gateway. status=%d", xcontext.Status(ctx))
	}

	var (
		w   life.Workable
		oid int64
		err error
	)

	if oid, err = xcontext.OID(ctx); err != nil {
		return err
	}

	replyFunc := func(p proto.Message) error {
		msg, ok := p.(*intrav1.TunnelResponse)
		if !ok {
			return errors.Wrapf(err, "intrav1.TunnelResponse proto type conversion failed")
		}
		logging.Reply(ctx, s.log, msg, logging.DefaultFilter)
		if err0 := stream.Send(msg); err0 != nil {
			return err0
		}
		return nil
	}

	if w, err = s.mgr.Worker(ctx, oid, core.NewReplier(replyFunc), life.NewBroadcaster(s.mgr.Pusher())); err != nil {
		return err
	}
	return s.run(ctx, w, stream)
}

func (s *TunnelService) run(ctx context.Context, w life.Workable, stream intrav1.TunnelService_TunnelServer) error {
	var err error

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		<-ctx.Done()
		return ctx.Err()
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			for {
				var in *intrav1.TunnelRequest
				if in, err = stream.Recv(); err != nil {
					return err
				}
				logging.Req(ctx, s.log, in, logging.DefaultFilter)
				if err = w.ProductFuncEvent(func(wctx life.Context) error {
					return s.handle(wctx.(core.Context), in)
				}); err != nil {
					return err
				}
			}
		})
	})

	if err = eg.Wait(); err != nil {
		id := w.ID()
		s.log.WithContext(ctx).Debugf("tunnel is stopping. uid=%d signal=%v", id, err.Error())
		w.Stop(ctx)
		s.log.WithContext(ctx).Debugf("tunnel stopped. uid=%d", id)
		return nil
	}
	return nil
}

func (s *TunnelService) handle(wctx core.Context, in *intrav1.TunnelRequest) error {
	out, err := handler.PlayerHandle(wctx, s.svc, in)
	if err != nil {
		s.log.WithContext(wctx).Errorf("handle failed. %d-%d uid=%d %+v", in.Mod, in.Seq, wctx.UID(), err)
	}
	if out == nil {
		sc := &climsg.SCServerUnexpectedErr{
			Mod: in.Mod,
			Seq: in.Seq,
		}
		if out, err = handler.NewPlayerResponse(int32(climod.ModuleID_System), int32(cliseq.SystemSeq_ServerUnexpectedErr), in.Obj, sc); err != nil {
			return errors.Wrapf(err, "TunnelResponse marshal failed. out=SCServerUnexpectedErr uid=%d", wctx.UID())
		}
	}
	return wctx.ReplyBytes(climod.ModuleID(in.Mod), in.Seq, in.Obj, out)
}
