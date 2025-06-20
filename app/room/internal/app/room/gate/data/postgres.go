package data

import (
	"context"

	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	"github.com/go-pantheon/roma/app/room/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
)

var _ domain.RoomRepo = (*roomPostgresRepo)(nil)

type roomPostgresRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewRoomPostgresRepo(data *data.Data, logger log.Logger) (domain.RoomRepo, error) {
	r := &roomPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "room/room/gate/data")),
	}

	return r, nil
}

func (r *roomPostgresRepo) Create(ctx context.Context, p *dbv1.RoomProto, ctime time.Time) (err error) {
	return errors.New("not implemented")
}

func (r *roomPostgresRepo) QueryByID(ctx context.Context, id int64, p *dbv1.RoomProto) error {
	return errors.New("not implemented")
}

func (r *roomPostgresRepo) UpdateByID(ctx context.Context, id int64, user *dbv1.RoomProto) error {
	return errors.New("not implemented")
}

func (r *roomPostgresRepo) Exist(ctx context.Context, id int64) (bool, error) {
	return false, errors.New("not implemented")
}

func (r *roomPostgresRepo) IncVersion(ctx context.Context, id int64, newVersion int64) error {
	return errors.New("not implemented")
}
