package server

import (
	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/intra/filter"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/server/registry"
)

func NewHTTPServer(
	c *conf.Server, logger log.Logger, filter *filter.HttpFilter, svcRg *registry.ServiceRegistrars, gateRg *registry.GateRegistrars) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
				logging.Server(logger),
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
	return svr
}
