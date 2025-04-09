package registry

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/service"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/gamedata/v1"
)

type GamedataRegistrar struct {
	gamedata *service.GamedataAdmin
}

func NewGamedataRegistrar(
	gamedata *service.GamedataAdmin,
) *GamedataRegistrar {
	return &GamedataRegistrar{
		gamedata: gamedata,
	}
}

func (r *GamedataRegistrar) GrpcRegister(s *grpc.Server) {
	adminv1.RegisterGamedataAdminServer(s, r.gamedata)
	log.Infof("Register gamedata admin gRPC service")
}

func (r *GamedataRegistrar) HttpRegister(s *http.Server) {
	adminv1.RegisterGamedataAdminHTTPServer(s, r.gamedata)
	log.Infof("Register gamedata admin HTTP service")
}
