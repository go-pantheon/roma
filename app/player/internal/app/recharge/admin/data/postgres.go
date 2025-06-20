package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

var _ domain.OrderRepo = (*orderPostgresRepo)(nil)

type orderPostgresRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewOrderPostgresRepo(data *data.Data, logger log.Logger) (domain.OrderRepo, error) {
	r := &orderPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/admin/data/order")),
	}

	return r, nil
}

func (r *orderPostgresRepo) GetByID(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error) {
	return nil, errors.New("not implemented")
}

func (r *orderPostgresRepo) GetList(ctx context.Context, index, limit int64, cond *dbv1.OrderProto) ([]*dbv1.OrderProto, int64, error) {
	return nil, 0, errors.New("not implemented")
}

func (r *orderPostgresRepo) UpdateAckStateByID(ctx context.Context, store pkg.Store, transId string, state dbv1.OrderAckState) error {
	return errors.New("not implemented")
}
