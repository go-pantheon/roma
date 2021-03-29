package biz

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/domain"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/core"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/room/admin/room/v1"
)

type RoomUseCase struct {
	log *log.Helper
	mgr *core.Manager
	do  *domain.RoomDomain
}

func NewRoomUseCase(mgr *core.Manager,
	do *domain.RoomDomain,
	logger log.Logger,
) *RoomUseCase {
	uc := &RoomUseCase{
		log: log.NewHelper(log.With(logger, "module", "room/admin/biz/room")),
		mgr: mgr,
		do:  do,
	}

	return uc
}

func (uc *RoomUseCase) GetById(ctx context.Context, id int64) (g *adminv1.RoomProto, err error) {
	return uc.do.Load(ctx, id)
}
