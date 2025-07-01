package registry

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

func NewRoomRegistrar(svc climsg.RoomServiceServer) *RoomRegistrar {
	return &RoomRegistrar{
		svc: svc,
	}
}

type RoomRegistrar struct {
	svc climsg.RoomServiceServer
}

func (r *RoomRegistrar) GrpcRegister(s *grpc.Server) {
	climsg.RegisterRoomServiceServer(s, r.svc)
	log.Infof("Register room gRPC service")
}

func (r *RoomRegistrar) HttpRegister(s *http.Server) {
	climsg.RegisterRoomServiceHTTPServer(s, r.svc)
	log.Infof("Register room HTTP service")
}
