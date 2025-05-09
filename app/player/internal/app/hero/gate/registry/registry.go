package registry

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

func NewHeroRegistrar(svc climsg.HeroServiceServer) *HeroRegistrar {
	return &HeroRegistrar{
		svc: svc,
	}
}

type HeroRegistrar struct {
	svc climsg.HeroServiceServer
}

func (r *HeroRegistrar) GrpcRegister(s *grpc.Server) {
	climsg.RegisterHeroServiceServer(s, r.svc)
	log.Infof("Register hero tcp gRPC service")
}

func (r *HeroRegistrar) HttpRegister(s *http.Server) {
	climsg.RegisterHeroServiceHTTPServer(s, r.svc)
	log.Infof("Register hero tcp HTTP service")
}
