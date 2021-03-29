package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/room/admin/room/v1"
)

func NewRoomRegistrar(svc adminv1.RoomAdminServer) *RoomRegistrar {
	return &RoomRegistrar{
		svc: svc,
	}
}

type RoomRegistrar struct {
	svc adminv1.RoomAdminServer
}

func (r *RoomRegistrar) GrpcRegister(s *grpc.Server) {
	adminv1.RegisterRoomAdminServer(s, r.svc)
	log.Infof("Register room admin gRPC service")
}

func (r *RoomRegistrar) HttpRegister(s *http.Server) {
	adminv1.RegisterRoomAdminHTTPServer(s, r.svc)
	log.Infof("Register room admin HTTP service")
}
