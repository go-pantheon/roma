package server

import (
	"math"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-pantheon/fabrica-kit/metrics"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/intra/filter"
	"github.com/go-pantheon/roma/app/player/internal/server/registry"
	"google.golang.org/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server, logger log.Logger, filter *filter.GrpcFilter,
	svcRg *registry.ServiceRegistrars, adminRg *registry.AdminRegistrars,
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
	for _, r := range svcRg.Rgs {
		r.GrpcRegister(svr)
	}
	for _, r := range adminRg.Rgs {
		r.GrpcRegister(svr)
	}
	return svr
}
