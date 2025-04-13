package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/core"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

type RechargeUseCase struct {
	log *log.Helper
	mgr *core.Manager
	do  *domain.OrderDomain
}

func NewRechargeUseCase(mgr *core.Manager, do *domain.OrderDomain, logger log.Logger) *RechargeUseCase {
	uc := &RechargeUseCase{
		log: log.NewHelper(log.With(logger, "module", "player/recharge/admin/biz")),
		mgr: mgr,
		do:  do,
	}

	return uc
}

func (uc *RechargeUseCase) GetList(ctx context.Context, index, size int64, cond *dbv1.OrderProto) ([]*dbv1.OrderProto, int64, error) {
	return uc.do.GetList(ctx, index, size, cond)
}

func (uc *RechargeUseCase) GetById(ctx context.Context, store pkg.Store, transId string) (p *dbv1.OrderProto, err error) {
	return uc.do.GetByID(ctx, store, transId)
}

func (uc *RechargeUseCase) UpdateAckStateByID(ctx context.Context, store pkg.Store, transId string, state dbv1.OrderAckState) (err error) {
	return uc.do.UpdateAckStateByID(ctx, store, transId, state)
}
