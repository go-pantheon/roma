package service

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/recharge/v1"
)

type RechargeAdmin struct {
	adminv1.UnimplementedRechargeAdminServer

	log *log.Helper
	uc  *biz.RechargeUseCase
}

func NewRechargeAdmin(logger log.Logger, uc *biz.RechargeUseCase) adminv1.RechargeAdminServer {
	return &RechargeAdmin{
		log: log.NewHelper(log.With(logger, "module", "player/recharge/admin/service")),
		uc:  uc,
	}
}

func (s *RechargeAdmin) GetOrderList(ctx context.Context, req *adminv1.GetOrderListRequest) (*adminv1.GetOrderListResponse, error) {
	cond, start, limit, err := buildGetOrderListCond(req)
	if err != nil {
		return nil, err
	}

	protos, count, err := s.uc.GetList(ctx, start, limit, cond)
	if err != nil {
		return nil, err
	}

	reply := &adminv1.GetOrderListResponse{
		Orders: make([]*adminv1.OrderProto, 0, len(protos)),
		Total:  count,
	}

	for _, p := range protos {
		u, err := toOrderProto(p)
		if err != nil {
			s.log.WithContext(ctx).Errorf("proto to order proto failed. %+v", err)
			continue
		}

		reply.Orders = append(reply.Orders, u)
	}

	return reply, nil
}

func buildGetOrderListCond(req *adminv1.GetOrderListRequest) (cond *dbv1.OrderProto, start, limit int64, err error) {
	start, limit = profile.PageStartLimit(req.Page, req.PageSize)

	cond = &dbv1.OrderProto{}

	if req.Cond == nil {
		return nil, 0, 0, xerrors.APIParamInvalid("condition is empty")
	}

	if len(req.Cond.Store) > 0 {
		cond.Store = req.Cond.Store
	}

	if len(req.Cond.TransId) > 0 {
		cond.TransId = req.Cond.TransId
	}

	if req.Cond.Uid > 0 {
		cond.Uid = req.Cond.Uid
	}

	if req.Cond.Ack > 0 {
		cond.Ack = req.Cond.Ack
	}

	return cond, start, limit, nil
}

func (s *RechargeAdmin) GetOrderById(ctx context.Context, req *adminv1.GetOrderByIdRequest) (*adminv1.GetOrderByIdResponse, error) {
	if req.TransId == "" {
		return nil, xerrors.APIParamInvalid("transId is empty")
	}

	if req.Store == "" {
		return nil, xerrors.APIParamInvalid("store is empty")
	}

	store, err := pkg.StoreFromString(req.Store)
	if err != nil {
		return nil, xerrors.APIParamInvalid("invalid store").WithCause(err)
	}

	p, err := s.uc.GetById(ctx, store, req.TransId)
	if err != nil {
		return nil, err
	}

	u, err := toOrderProto(p)
	if err != nil {
		return nil, err
	}

	return &adminv1.GetOrderByIdResponse{
		Code:  http.StatusOK,
		Order: u,
	}, nil
}

func (s *RechargeAdmin) UpdateOrderAckStateById(ctx context.Context, req *adminv1.UpdateOrderAckStateByIdRequest) (*adminv1.UpdateOrderAckStateByIdResponse, error) {
	if req.TransId == "" {
		return nil, xerrors.APIParamInvalid("transId is empty")
	}

	if req.Store == "" {
		return nil, xerrors.APIParamInvalid("store is empty")
	}

	store, err := pkg.StoreFromString(req.Store)
	if err != nil {
		return nil, xerrors.APIParamInvalid("invalid store").WithCause(err)
	}

	if _, ok := dbv1.OrderAckState_name[req.Ack]; !ok {
		return nil, xerrors.APIParamInvalid("ack state not exists")
	}

	if req.Ack != int32(dbv1.OrderAckState_ORDER_ACK_STATE_GM_SUCCEEDED) && req.Ack != int32(dbv1.OrderAckState_ORDER_ACK_STATE_GM_CANCELED) {
		return nil, xerrors.APIParamInvalid("invalid ack state")
	}

	if err := s.uc.UpdateAckStateByID(ctx, store, req.TransId, dbv1.OrderAckState(req.Ack)); err != nil {
		return nil, err
	}

	return &adminv1.UpdateOrderAckStateByIdResponse{
		Code: http.StatusOK,
	}, nil
}

func toOrderProto(p *dbv1.OrderProto) (*adminv1.OrderProto, error) {
	p.Info.Token = ""

	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, xerrors.APICodecFailed("json marshal failed").WithCause(err)
	}

	u := &adminv1.OrderProto{
		Store:       p.Store,
		TransId:     p.TransId,
		Uid:         p.Uid,
		Ack:         p.Ack,
		ProductId:   p.Info.ProductId,
		PurchasedAt: p.Info.PurchasedAt,
		AckAt:       p.AckAt,
		Detail:      string(bytes),
	}

	return u, nil
}
