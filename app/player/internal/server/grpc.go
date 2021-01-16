package server

import (
	"math"

	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/middleware"
	"github.com/go-kratos/kratos/middleware/logging"
	"github.com/go-kratos/kratos/middleware/metadata"
	"github.com/go-kratos/kratos/middleware/recovery"
	"github.com/go-kratos/kratos/middleware/tracing"
	kgrpc "github.com/go-kratos/kratos/transport/grpc"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/intra/filter"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/server/registry"
	"github.com/vulcan-frame/vulcan-kit/metrics"
	"google.golang.org/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server, logger log.Logger, filter *filter.GrpcFilter, rg *registry.ServiceRegistrars,
) *kgrpc.Server {
	var opts = []kgrpc.ServerOption{
		kgrpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
				tracing.Server(),
				metrics.Server(),
				logging.Server(logger),
				filter.Server(),
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
	opts = append(opts, kgrpc.Options(
		grpc.InitialConnWindowSize(1<<30),
		grpc.InitialWindowSize(1<<30),
		grpc.MaxConcurrentStreams(math.MaxInt32),
	))

	svr := kgrpc.NewServer(opts...)
	for _, r := range rg.Rgs {
		r.GrpcRegister(svr)
	}
	return svr
}
