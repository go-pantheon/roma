package filter

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/roma/app/room/internal/core"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/go-pantheon/roma/pkg/universe/life"
	unimd "github.com/go-pantheon/roma/pkg/universe/middleware/dev"
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
