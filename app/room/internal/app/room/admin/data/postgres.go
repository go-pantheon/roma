package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/domain"
	"github.com/go-pantheon/roma/app/room/internal/data"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/room/admin/room/v1"
)

var _ domain.RoomRepo = (*roomPostgresRepo)(nil)

type roomPostgresRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewRoomPostgresRepo(data *data.Data, logger log.Logger) (domain.RoomRepo, error) {
	r := &roomPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/admin/data/order")),
	}

	return r, nil
}

func (r *roomPostgresRepo) GetByID(ctx context.Context, id int64) (*adminv1.RoomProto, error) {
	return nil, errors.New("not implemented")
}
