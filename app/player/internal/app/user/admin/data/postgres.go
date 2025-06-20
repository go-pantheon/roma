package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

var _ domain.UserRepo = (*userPostgresRepo)(nil)

type userPostgresRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewUserPostgresRepo(data *data.Data, logger log.Logger) (domain.UserRepo, error) {
	r := &userPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/user/gate/data")),
	}

	return r, nil
}

func (r *userPostgresRepo) GetByID(ctx context.Context, uid int64) (*dbv1.UserProto, error) {
	return nil, errors.New("not implemented")
}

func (r *userPostgresRepo) GetList(ctx context.Context, start, limit int64, cond *dbv1.UserProto) (result []*dbv1.UserProto, count int64, err error) {
	return nil, 0, errors.New("not implemented")
}

func (r *userPostgresRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	return errors.New("not implemented")
}
