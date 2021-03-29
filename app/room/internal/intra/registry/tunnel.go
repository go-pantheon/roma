package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	intrav1 "github.com/vulcan-frame/vulcan-game/gen/api/server/room/intra/v1"
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
	log.Infof("Register infra HTTP service")
}
