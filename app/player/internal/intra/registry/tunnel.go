package registry

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/player/intra/v1"
)

type IntraRegistrar struct {
	tunnel intrav1.TunnelServiceServer
}

func NewIntraRegistrar(tunnel intrav1.TunnelServiceServer) *IntraRegistrar {
	return &IntraRegistrar{
		tunnel: tunnel,
	}
}

func (r *IntraRegistrar) GrpcRegister(s *grpc.Server) {
	intrav1.RegisterTunnelServiceServer(s, r.tunnel)
	log.Infof("Register infra gRPC service")
}

func (r *IntraRegistrar) HttpRegister(s *http.Server) {
	log.Infof("Ignore register infra HTTP service")
}
