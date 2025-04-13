package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/room/admin/room/v1"
)

type RoomRepo interface {
	GetByID(ctx context.Context, id int64) (*adminv1.RoomProto, error)
}

type RoomDomain struct {
	log  *log.Helper
	repo RoomRepo
}

func NewRoomDomain(pr RoomRepo, logger log.Logger) *RoomDomain {
	return &RoomDomain{
		repo: pr,
		log:  log.NewHelper(log.With(logger, "module", "room/admin/domain/room"))}
}

func (do *RoomDomain) Load(ctx context.Context, id int64) (u *adminv1.RoomProto, err error) {
	return do.repo.GetByID(ctx, id)
}
