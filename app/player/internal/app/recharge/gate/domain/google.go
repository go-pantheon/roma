package domain

import (
	"context"
	"strconv"
	"time"

	"github.com/awa/go-iap/playstore"
	"github.com/awa/go-iap/playstore/mocks"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/rechargeerrs"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/go-pantheon/roma/app/player/internal/core"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	gomock "go.uber.org/mock/gomock"
	"google.golang.org/api/androidpublisher/v3"
	"google.golang.org/protobuf/proto"
)

type googleReceipt struct {
	Json      string `json:"json"`
	Signature string `json:"signature"`
}

type googleReceiptJson struct {
	OrderId          string `json:"orderId"`
	PackageName      string `json:"packageName"`
	ProductId        string `json:"productId"`
	PurchaseTime     int64  `json:"purchaseTime"`
	PurchaseState    int    `json:"purchaseState"`
	DeveloperPayload string `json:"developerPayload"`
	PurchaseToken    string `json:"purchaseToken"`
}

var _ Verifiable = (*googleCli)(nil)

type googleCli struct {
	store pkg.Store
	log   *log.Helper
	cli   playstore.IABProduct
	conf  *conf.Recharge_Google
	repo  OrderRepo
}

func newGoogleCli(logger log.Logger, label *conf.Label, c *conf.Recharge_Google, repo OrderRepo) (do *googleCli, err error) {
	do = &googleCli{
		store: pkg.StoreGoogle,
		log:   log.NewHelper(log.With(logger, "module", "player/recharge/gate/domain/google")),
		conf:  proto.Clone(c).(*conf.Recharge_Google),
		repo:  repo,
	}

	if profile.IsDevStr(label.Profile) {
		do.cli = mocks.NewMockIABProduct(gomock.NewController(nil))

		do.log.Debugf("use mock google recharge client")

		return do, nil
	}

	do.cli, err = playstore.New([]byte(c.Json))
	if err != nil {
		return nil, errors.Wrapf(err, "create google api failed")
	}

	return do, nil
}

func (do *googleCli) verifyReceipt(ctx core.Context, receipt []byte, productId int64) (reset *ResetOrderInfo, err error) {
	var json *googleReceiptJson

	if json, err = do.buildJson(ctx, receipt, productId); err != nil {
		return nil, err
	}

	ctx0, cancelFunc := context.WithTimeout(ctx, ttl)
	defer cancelFunc()

	resp, err := do.cli.VerifyProduct(ctx0, json.PackageName, json.ProductId, json.PurchaseToken)
	if err != nil {
		return nil, errors.Wrapf(rechargeerrs.ErrApiVerify, "%s", err.Error())
	}

	// check api resp. PurchaseState: 0. Purchased 1. Canceled 2. Pending
	if resp.PurchaseState != 0 {
		return nil, errors.Wrapf(rechargeerrs.ErrNotPurchased, "purchase_state=%d", resp.PurchaseState)
	}

	// api verify succeeded, check order is existed
	if err = checkOrder(ctx, do.repo, do.store, json.OrderId); err != nil {
		return nil, err
	}

	order := do.createOrderProto(ctx.UID(), json, resp, ctx.Now())
	err = do.repo.Create(ctx, order)

	reset = &ResetOrderInfo{
		Store:    do.store,
		TransIds: []string{order.TransId},
	}

	return reset, nil
}

func (do *googleCli) buildJson(ctx context.Context, receipt []byte, productId int64) (json *googleReceiptJson, err error) {
	gr := &googleReceipt{}

	if err = jsoniter.Unmarshal(receipt, gr); err != nil {
		do.log.WithContext(ctx).Errorf("unmarshal receipt failed. %s", err.Error())
		return nil, err
	}

	isValid, err := playstore.VerifySignature(do.conf.PubKey, []byte(gr.Json), gr.Signature)
	if err != nil {
		return nil, errors.Wrapf(rechargeerrs.ErrSignature, "%s", err.Error())
	}

	if !isValid {
		return nil, errors.Wrapf(rechargeerrs.ErrSignature, "signature=%s", gr.Signature)
	}

	json = &googleReceiptJson{}

	if err = jsoniter.UnmarshalFromString(gr.Json, json); err != nil {
		return nil, errors.Wrapf(rechargeerrs.ErrUnmarshal, "%s", err.Error())
	}

	if json.ProductId != strconv.FormatInt(productId, 10) {
		return nil, errors.Wrapf(rechargeerrs.ErrProductId, "productId=%d resp=%s", productId, json.ProductId)
	}

	if json.PackageName != do.conf.PackageName {
		return nil, errors.Wrapf(rechargeerrs.ErrPackageName, "package_name=%s expect=%s", json.PackageName, do.conf.PackageName)
	}

	return json, nil
}

func (do *googleCli) createOrderProto(uid int64, receipt *googleReceiptJson, result *androidpublisher.ProductPurchase, ctime time.Time) *dbv1.OrderProto {
	p := &dbv1.OrderProto{
		Store: string(do.store),
		Uid:   uid,
		Info:  &dbv1.OrderInfoProto{},
	}

	p.Info.Token = receipt.PurchaseToken
	p.Info.ProductId = receipt.ProductId

	p.TransId = result.OrderId
	p.Info.PurchasedAt = result.PurchaseTimeMillis / microToMilli // microseconds to milliseconds

	p.Ack = int32(dbv1.OrderAckState_ORDER_ACK_STATE_SUCCEEDED)
	p.AckAt = ctime.Unix()

	return p
}
