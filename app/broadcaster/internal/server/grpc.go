package server

import (
	"math"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-pantheon/fabrica-kit/metrics"
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster/service"
	"github.com/go-pantheon/roma/app/broadcaster/internal/conf"
	v1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
	grpcgo "google.golang.org/grpc"
)

func NewGRPCServer(c *conf.Server, logger log.Logger, svc *service.BroadcasterService) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			middleware.Chain(
				recovery.Recovery(),
				metadata.Server(),
				tracing.Server(),
				metrics.Server(),
				logging.Server(logger),
			),
		),
	}

	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}

	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}

	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}

	opts = append(opts, grpc.Options(
		grpcgo.InitialConnWindowSize(1<<30),
		grpcgo.InitialWindowSize(1<<30),
		grpcgo.MaxConcurrentStreams(math.MaxInt32),
	))

	svr := grpc.NewServer(opts...)

	v1.RegisterPushServiceServer(svr, svc)

	return svr
}
