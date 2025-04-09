package registry

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

type DevRegistrar struct {
	svc climsg.DevServiceServer
}

func NewDevRegistrar(svc climsg.DevServiceServer) *DevRegistrar {
	return &DevRegistrar{
		svc: svc,
	}
}

func (r *DevRegistrar) GrpcRegister(s *grpc.Server) {
	climsg.RegisterDevServiceServer(s, r.svc)
	log.Infof("Register dev gRPC service")
}

func (r *DevRegistrar) HttpRegister(s *http.Server) {
	climsg.RegisterDevServiceHTTPServer(s, r.svc)
	log.Infof("Register dev HTTP service")
}
