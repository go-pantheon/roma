package service

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/errors"
	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/biz"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/room/admin/room/v1"
)

type RoomAdmin struct {
	adminv1.UnimplementedRoomAdminServer

	log *log.Helper
	uc  *biz.RoomUseCase
}

func NewRoomAdmin(logger log.Logger, uc *biz.RoomUseCase) adminv1.RoomAdminServer {
	return &RoomAdmin{
		log: log.NewHelper(log.With(logger, "module", "room/admin/service/room")),
		uc:  uc,
	}
}

func (s *RoomAdmin) GetById(ctx context.Context, req *adminv1.GetByIdRequest) (*adminv1.GetByIdResponse, error) {
	p, err := s.uc.GetById(ctx, req.Id)
	if err != nil {
		return nil, errors.BadRequest("查询失败", err.Error())
	}
	reply := &adminv1.GetByIdResponse{
		Code: http.StatusOK,
		Item: p,
	}
	return reply, nil
}
