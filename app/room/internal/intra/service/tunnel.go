package service

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/app/room/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/room/intra/v1"
	"github.com/go-pantheon/roma/gen/app/room/handler"
	"github.com/go-pantheon/roma/gen/app/room/service"
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
	svc *service.RoomServices
}

func NewTunnelService(logger log.Logger, mgr *core.Manager, svc *service.RoomServices) intrav1.TunnelServiceServer {
	return &TunnelService{
		log: log.NewHelper(log.With(logger, "module", "room/intra/tunnel")),
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

		if err := stream.Send(msg); err != nil {
			return errors.Wrapf(err, "intrav1.TunnelResponse send failed")
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
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.StopTriggered():
			return xsync.ErrStopByTrigger
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

	if err = eg.Wait(); err != nil {
		s.log.WithContext(ctx).Debugf("tunnel stopped. uid=%d signal=%v", w.ID(), err.Error())
		return nil
	}
	return nil
}

func (s *TunnelService) recv(w life.Workable, stream intrav1.TunnelService_TunnelServer) (err error) {
	var in *intrav1.TunnelRequest

	if in, err = stream.Recv(); err != nil {
		return errors.Wrapf(err, "intrav1.TunnelRequest recv failed")
	}

	if in == nil {
		return nil
	}

	return w.EmitFuncEvent(func(wctx life.Context) error {
		return s.handle(wctx.(core.Context), in)
	})
}

func (s *TunnelService) handle(wctx core.Context, in *intrav1.TunnelRequest) error {
	st := time.Now()
	logging.Req(wctx, s.log, in, logging.DefaultFilter)

	out, err := handler.RoomHandle(wctx, s.svc, in)
	if err != nil {
		// only log the handle error for keep the worker running
		s.log.WithContext(wctx).Errorf("room handle failed. uid=%d in=%d-%d obj=%d color=%d status=%d %+v", wctx.UID(), in.Mod, in.Seq, in.Obj, xcontext.Color(wctx), xcontext.Status(wctx), err)
	}

	if out == nil {
		if out, err = handler.NewRoomResponse(
			int32(climod.ModuleID_System),
			int32(cliseq.SystemSeq_ServerUnexpectedErr),
			in.Obj,
			&climsg.SCServerUnexpectedErr{
				Mod: in.Mod,
				Seq: in.Seq,
			},
		); err != nil {
			return errors.Wrapf(err, "TunnelResponse marshal failed. out=SCServerUnexpectedErr uid=%d", wctx.UID())
		}
	}

	logging.Reply(wctx, s.log, wctx.UID(), in, out, time.Since(st), logging.DefaultFilter)

	return wctx.ReplyBytes(climod.ModuleID(in.Mod), in.Seq, in.Obj, out)
}
