package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-pantheon/fabrica-kit/metrics"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/intra/filter"
	"github.com/go-pantheon/roma/app/player/internal/server/registry"
	devmd "github.com/go-pantheon/roma/pkg/universe/middleware/dev"
)

func NewHTTPServer(
	c *conf.Server, logger log.Logger, filter *filter.HttpFilter, 
	svcRg *registry.ServiceRegistrars, gateRg *registry.GateRegistrars, adminRg *registry.AdminRegistrars,
) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
				tracing.Server(),
				metrics.Server(),
				devmd.Server(),
				logging.Server(logger),
				filter.Server(),
			)),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	svr := http.NewServer(opts...)
	for _, r := range gateRg.Rgs {
		r.HttpRegister(svr)
	}
	for _, r := range svcRg.Rgs {
		r.HttpRegister(svr)
	}
	for _, r := range adminRg.Rgs {
		r.HttpRegister(svr)
	}
	return svr
}
