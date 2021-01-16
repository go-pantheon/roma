package server

import (
	"math"

	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/middleware"
	"github.com/go-kratos/kratos/middleware/metadata"
	"github.com/go-kratos/kratos/middleware/recovery"
	kgrpc "github.com/go-kratos/kratos/transport/grpc"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/server/registry"
	"google.golang.org/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server, logger log.Logger, rg *registry.ServiceRegistrars,
) *kgrpc.Server {
	var opts = []kgrpc.ServerOption{
		kgrpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
			),
		),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, kgrpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, kgrpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, kgrpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	svr := kgrpc.NewServer(opts...)
	for _, r := range rg.Rgs {
		r.GrpcRegister(svr)
	}
	return svr
}
