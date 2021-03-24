package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/recharge/v1"
)

type RechargeRegistrar struct {
	recharge adminv1.RechargeAdminServer
}

func NewRechargeRegistrar(
	recharge adminv1.RechargeAdminServer,
) *RechargeRegistrar {
	return &RechargeRegistrar{
		recharge: recharge,
	}
}

func (r *RechargeRegistrar) GrpcRegister(s *grpc.Server) {
	adminv1.RegisterRechargeAdminServer(s, r.recharge)
	log.Infof("Register recharge admin gRPC service")
}

func (r *RechargeRegistrar) HttpRegister(s *http.Server) {
	adminv1.RegisterRechargeAdminHTTPServer(s, r.recharge)
	log.Infof("Register recharge admin HTTP service")
}
