package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/user/v1"
)

type UserRegistrar struct {
	user adminv1.UserAdminServer
}

func NewUserRegistrar(
	user adminv1.UserAdminServer,
) *UserRegistrar {
	return &UserRegistrar{
		user: user,
	}
}

func (r *UserRegistrar) GrpcRegister(s *grpc.Server) {
	adminv1.RegisterUserAdminServer(s, r.user)
	log.Infof("Register user admin gRPC service")
}

func (r *UserRegistrar) HttpRegister(s *http.Server) {
	adminv1.RegisterUserAdminHTTPServer(s, r.user)
	log.Infof("Register user admin HTTP service")
}
