package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/rechargeerrs"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/pkg/errors"
)

type Verifiable interface {
	verifyReceipt(ctx core.Context, receipt []byte, productId int64) (reset *ResetOrderInfo, err error) // verify receipt and create order
}

type OrderRepo interface {
	Create(ctx context.Context, order *dbv1.OrderProto) (err error)
	GetByTransId(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error)
	UpdateAckState(ctx context.Context, store pkg.Store, transId string, ack dbv1.OrderAckState) error
}

type Receipts struct {
	Store         pkg.Store `json:"Store"`
	TransactionID string    `json:"TransactionID"`
	Payload       string    `json:"Payload"`
}

type ResetOrderInfo struct {
	Store    pkg.Store
	TransIds []string
}

type RechargeDomain struct {
	log    *log.Helper
	repo   OrderRepo
	google *googleCli
	apple  *appleCli
}

func NewRechargeDomain(logger log.Logger, label *conf.Label, c *conf.Recharge, repo OrderRepo) (do *RechargeDomain, err error) {
	do = &RechargeDomain{
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/gate/domain/recharge")),
		repo: repo,
	}

	do.apple = newAppleCli(logger, repo, c.Apple)
	do.google, err = newGoogleCli(logger, label, c.Google, repo)

	if err != nil {
		return
	}

	return do, nil
}

func (do *RechargeDomain) VerifyRecharge(ctx core.Context, arg *climsg.RechargeParamProto, productId int64) (reset *ResetOrderInfo, err error) {
	if arg == nil {
		return nil, errors.Errorf("recharge arg is nil")
	}

	var cli Verifiable

	switch pkg.Store(arg.Store) {
	case pkg.StoreGoogle:
		cli = do.google
	case pkg.StoreApple:
		cli = do.apple
	default:
		return nil, errors.Wrapf(rechargeerrs.ErrRechargeType, "store=%s", arg.Store)
	}

	if profile.IsDev() {
		if err = do.AddUserRecharge(ctx, productId); err != nil {
			return nil, err
		}
		// return if dev
		return &ResetOrderInfo{
			Store: pkg.Store(arg.Store),
		}, nil
	}

	reset, err = cli.verifyReceipt(ctx, arg.Payload, productId)
	if err != nil {
		return nil, err
	}

	if err = do.AddUserRecharge(ctx, productId); err != nil {
		return nil, err
	}

	return reset, nil
}

func (do *RechargeDomain) ResetOrderAck(ctx core.Context, reset *ResetOrderInfo) {
	switch reset.Store {
	case pkg.StoreGoogle, pkg.StoreApple:
	default:
		do.log.WithContext(ctx).Errorf("reset order ack failed, unknown store. store=%s", reset.Store)
		return
	}

	for _, transId := range reset.TransIds {
		if err := do.repo.UpdateAckState(ctx, reset.Store, transId, dbv1.OrderAckState_ORDER_ACK_STATE_PENDING); err != nil {
			do.log.WithContext(ctx).Errorf("reset order ack failed. transId=%s, store=%s, err=%s", transId, reset.Store, err.Error())
		}
	}
}

func checkOrder(ctx context.Context, repo OrderRepo, store pkg.Store, transId string) error {
	order, err := repo.GetByTransId(ctx, store, transId)
	if err != nil {
		if errors.Is(err, xerrors.ErrDBRecordNotFound) {
			return nil
		}

		return err
	}

	if order == nil {
		return nil
	}

	return errors.Wrapf(rechargeerrs.ErrExisted, "store=%s transId=%s", store, transId)
}

func (do *RechargeDomain) AddUserRecharge(ctx core.Context, productId int64) error {
	d := gamedata.GetRechargeProductData(productId)
	if d == nil {
		return errors.Wrapf(rechargeerrs.ErrProductId, "productId=%d", productId)
	}

	if err := ctx.User().Recharge().AddRecharge(int64(d.Price)); err != nil {
		return errors.WithMessagef(err, "productId=%d price=%d", productId, d.Price)
	}

	ctx.Changed(object.ModuleKey)

	return nil
}
