package domain

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/awa/go-iap/appstore"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/rechargeerrs"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/core"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

const (
	ttl          = 10 * time.Second
	microToMilli = 1000
)

var _ Verifiable = (*appleCli)(nil)

type appleCli struct {
	store pkg.Store
	log   *log.Helper
	cli   appstore.IAPClient
	conf  *conf.Recharge_Apple
	repo  OrderRepo
}

func newAppleCli(logger log.Logger, repo OrderRepo, c *conf.Recharge_Apple) (do *appleCli) {
	do = &appleCli{
		store: pkg.StoreApple,
		log:   log.NewHelper(log.With(logger, "module", "player/recharge/gate/domain/apple")),
		cli:   appstore.New(),
		conf:  proto.Clone(c).(*conf.Recharge_Apple),
		repo:  repo,
	}
	return do
}

func (do *appleCli) verifyReceipt(ctx core.Context, receipt []byte, productId int64) (reset *ResetOrderInfo, err error) {
	ctx0, cancelFunc := context.WithTimeout(ctx, ttl)
	defer cancelFunc()

	var (
		req = appstore.IAPRequest{
			ReceiptData: string(receipt),
		}
		resp = &appstore.IAPResponse{}
	)

	if err = do.cli.Verify(ctx0, req, resp); err != nil {
		err = errors.Wrapf(rechargeerrs.ErrApiVerify, "%s", err.Error())
		return
	}
	if resp.Status != 0 {
		err = errors.Wrapf(rechargeerrs.ErrApiVerify, "status=%d", resp.Status)
		return
	}

	transIds, err := do.createOrder(ctx, req.ReceiptData, &resp.Receipt, resp.LatestReceiptInfo, productId)
	if err != nil {
		return
	}

	reset = &ResetOrderInfo{
		Store:    do.store,
		TransIds: transIds,
	}
	return
}

func (do *appleCli) createOrder(ctx core.Context, token string, receipt *appstore.Receipt, latestReceiptInfo []appstore.InApp, productId int64) (transIds []string, err error) {
	if receipt == nil {
		err = errors.Wrapf(rechargeerrs.ErrApiVerify, "receipt is nil")
		return
	}

	if receipt.BundleID != do.conf.BundleId {
		err = errors.Wrapf(rechargeerrs.ErrPackageName, "bundleId:%s", receipt.BundleID)
		return
	}

	for _, inApp := range receipt.InApp {
		if transId, err0 := do.createOrder0(ctx, productId, token, inApp); err0 != nil {
			do.log.WithContext(ctx).Errorf("verify receipt.inApp failed. %+v", err0)
		} else {
			if len(transId) > 0 {
				transIds = append(transIds, transId)
			}
		}
	}
	for _, inApp := range latestReceiptInfo {
		if transId, err0 := do.createOrder0(ctx, productId, token, inApp); err0 != nil {
			do.log.WithContext(ctx).Errorf("verify latestReceiptInfo failed. %+v", err0)
			continue
		} else {
			if len(transId) > 0 {
				transIds = append(transIds, transId)
			}
		}
	}

	if len(transIds) == 0 {
		err = errors.Wrapf(rechargeerrs.ErrProductId, "no valid transaction")
		return
	}
	return
}

func (do *appleCli) createOrder0(ctx core.Context, productId int64, token string, inapp appstore.InApp) (transId string, err error) {
	if inapp.ProductID != fmt.Sprintf("%s.%d", do.conf.BundleId, productId) {
		err = errors.Wrapf(rechargeerrs.ErrProductId, "ProductID=%s expect=%d", inapp.ProductID, productId)
		return
	}
	if inapp.InAppOwnershipType != "PURCHASED" {
		err = errors.Wrapf(rechargeerrs.ErrNotPurchased, "InAppOwnershipType=%s", inapp.InAppOwnershipType)
		return
	}
	// api verify succeeded, check order is existed
	if err = checkOrder(ctx, do.repo, do.store, inapp.TransactionID); err != nil {
		return
	}

	order := do.createOrderProto(ctx, ctx.UID(), token, &inapp, ctx.Now())
	if err = do.repo.Create(ctx, order); err != nil {
		return
	}

	transId = order.TransId
	return
}

func (do *appleCli) createOrderProto(_ context.Context, uid int64, token string, inApp *appstore.InApp, ctime time.Time) *dbv1.OrderProto {
	p := &dbv1.OrderProto{
		Store: string(do.store),
		Uid:   uid,
		Info:  &dbv1.OrderInfoProto{},
	}

	p.Info.Token = token
	p.Info.ProductId = inApp.ProductID
	p.TransId = inApp.TransactionID

	purchaseTime, err := strconv.ParseInt(inApp.PurchaseDateMS, 10, 64)
	if err != nil {
		p.Info.PurchasedAt = ctime.Unix()
	} else {
		p.Info.PurchasedAt = purchaseTime / microToMilli // microseconds to milliseconds
	}

	p.Ack = int32(dbv1.OrderAckState_ORDER_ACK_STATE_SUCCEEDED)
	p.AckAt = ctime.Unix()
	return p
}
