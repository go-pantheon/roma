package filter

import (
	"context"

	"github.com/go-kratos/kratos/middleware"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/core"
	"github.com/vulcan-frame/vulcan-game/pkg/errs"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/life"
	unimd "github.com/vulcan-frame/vulcan-game/pkg/universe/middleware/dev"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
)

type HttpFilter struct {
	mgr *core.Manager
}

func NewHttpFilter(mgr *core.Manager) *HttpFilter {
	md := &HttpFilter{
		mgr: mgr,
	}
	return md
}

func (md *HttpFilter) Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if !profile.IsDev() {
				return nil, errs.ErrProfileIllegal
			}

			ctx = unimd.TransformContext(ctx)
			oid, err := xcontext.OID(ctx)
			if err != nil {
				return nil, err
			}

			err = md.mgr.ExecuteAppEvent(ctx, oid, func(wctx life.Context) (err error) {
				reply, err = handler(wctx, req)
				return
			})
			return reply, err
		}
	}
}
