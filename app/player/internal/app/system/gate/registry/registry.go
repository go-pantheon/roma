package registry

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

func NewSystemRegistrar(svc climsg.SystemServiceServer) *SystemRegistrar {
	return &SystemRegistrar{
		svc: svc,
	}
}

type SystemRegistrar struct {
	svc climsg.SystemServiceServer
}

func (r *SystemRegistrar) GrpcRegister(s *grpc.Server) {
	climsg.RegisterSystemServiceServer(s, r.svc)
	log.Infof("Register system tcp gRPC service")
}

func (r *SystemRegistrar) HttpRegister(s *http.Server) {
	climsg.RegisterSystemServiceHTTPServer(s, r.svc)
	log.Infof("Register system tcp HTTP service")
}
