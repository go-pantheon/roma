package service

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
)

type RoomService struct {
	climsg.UnimplementedRoomServiceServer

	log *log.Helper
	uc  *biz.RoomUseCase
}

func NewRoomService(logger log.Logger, uc *biz.RoomUseCase) climsg.RoomServiceServer {
	return &RoomService{
		log: log.NewHelper(log.With(logger, "module", "room/room/gate/service")),
		uc:  uc,
	}
}

func (s *RoomService) RoomList(ctx context.Context, cs *climsg.CSRoomList) (sc *climsg.SCRoomList, err error) {
	return s.uc.RoomList(ctx.(core.Context), cs)
}

func (s *RoomService) CreateRoom(ctx context.Context, cs *climsg.CSCreateRoom) (sc *climsg.SCCreateRoom, err error) {
	return s.uc.CreateRoom(ctx.(core.Context), cs)
}
