package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

type OrderRepo interface {
	GetByID(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error)
	GetList(ctx context.Context, index, limit int64, cond *dbv1.OrderProto) ([]*dbv1.OrderProto, int64, error)
	UpdateAckStateByID(ctx context.Context, store pkg.Store, transId string, state dbv1.OrderAckState) error
}

type OrderDomain struct {
	log  *log.Helper
	repo OrderRepo
}

func NewOrderDomain(pr OrderRepo, logger log.Logger) *OrderDomain {
	return &OrderDomain{
		repo: pr,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/admin/domain/order")),
	}
}

func (do *OrderDomain) GetByID(ctx context.Context, store pkg.Store, transId string) (u *dbv1.OrderProto, err error) {
	return do.repo.GetByID(ctx, store, transId)
}

func (do *OrderDomain) GetList(ctx context.Context, index, size int64, cond *dbv1.OrderProto) ([]*dbv1.OrderProto, int64, error) {
	return do.repo.GetList(ctx, index, size, cond)
}

func (do *OrderDomain) UpdateAckStateByID(ctx context.Context, store pkg.Store, transId string, state dbv1.OrderAckState) (err error) {
	return do.repo.UpdateAckStateByID(ctx, store, transId, state)
}
