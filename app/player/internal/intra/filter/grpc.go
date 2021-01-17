package filter

import (
	"context"

	"github.com/go-kratos/kratos/middleware"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
)

type GrpcFilter struct {
	mgr *core.Manager
}

func NewGrpcFilter(mgr *core.Manager) *GrpcFilter {
	md := &GrpcFilter{mgr: mgr}
	return md
}

func (md *GrpcFilter) Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if !life.IsInnerContext(ctx) {
				return handler(ctx, req)
			}

			oid, err := xcontext.OID(ctx)
			if err != nil {
				return nil, err
			}

			err = md.mgr.ExecuteAppEvent(ctx, oid,
				func(wctx life.Context) (err error) {
					reply, err = handler(wctx, req)
					return
				})
			return reply, err
		}
	}
}
