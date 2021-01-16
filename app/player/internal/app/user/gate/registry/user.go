package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
)

type UserRegistrar struct {
	svc climsg.UserServiceServer
}

func NewUserRegistrar(svc climsg.UserServiceServer) *UserRegistrar {
	return &UserRegistrar{
		svc: svc,
	}
}

func (r *UserRegistrar) GrpcRegister(s *grpc.Server) {
	climsg.RegisterUserServiceServer(s, r.svc)
	log.Infof("Register user tcp gRPC service")
}

func (r *UserRegistrar) HttpRegister(s *http.Server) {
	climsg.RegisterUserServiceHTTPServer(s, r.svc)
	log.Infof("Register user tcp HTTP service")
}
