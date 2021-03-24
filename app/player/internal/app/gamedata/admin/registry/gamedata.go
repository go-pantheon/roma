package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/gamedata/admin/service"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/gamedata/v1"
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
