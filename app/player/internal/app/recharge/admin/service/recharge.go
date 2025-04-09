package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kratos/kratos/v2/errors"
	kerrors "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/recharge/v1"
	"github.com/go-pantheon/roma/pkg/util/maths/i32"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	maxPageSize     = 200
	defaultPageSize = 10
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
	cond, page, pageSize, err := buildGetOrderListCond(req)
	if err != nil {
		return nil, err
	}

	protos, count, err := s.uc.GetList(ctx, i32.Max(page-1, 0)*pageSize, pageSize, cond)
	if err != nil {
		return nil, kerrors.BadRequest(xerrors.ErrAdminQueryFailedReason, err.Error())
	}

	reply := &adminv1.GetOrderListResponse{
		Code:   http.StatusOK,
		Orders: make([]*adminv1.OrderProto, 0, len(protos)),
		Total:  uint32(count),
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

func buildGetOrderListCond(req *adminv1.GetOrderListRequest) (cond *dbv1.OrderProto, page, pageSize int32, err error) {
	if req.PageSize > maxPageSize {
		err = kerrors.BadRequest(xerrors.ErrAdminParamReason, fmt.Sprintf("page size must be less than %d", maxPageSize))
		return
	}

	if page = req.Page; page <= 0 {
		page = 1
	}
	if pageSize = req.PageSize; pageSize <= 0 {
		pageSize = defaultPageSize
	}

	cond = &dbv1.OrderProto{}
	if req.Cond == nil {
		err = kerrors.BadRequest(xerrors.ErrAdminConditionReason, "condition is empty")
		return
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
	return
}

func (s *RechargeAdmin) GetOrderById(ctx context.Context, req *adminv1.GetOrderByIdRequest) (*adminv1.GetOrderByIdResponse, error) {
	if req.TransId == "" {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "transId is empty")
	}
	if req.Store == "" {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "store is empty")
	}

	store, err := pkg.StoreFromString(req.Store)
	if err != nil {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "invalid store")
	}

	p, err := s.uc.GetById(ctx, store, req.TransId)
	if err != nil {
		return nil, kerrors.BadRequest(xerrors.ErrAdminQueryFailedReason, err.Error())
	}

	u, err := toOrderProto(p)
	if err != nil {
		return nil, kerrors.BadRequest(xerrors.ErrAdminQueryFailedReason, err.Error())
	}

	reply := &adminv1.GetOrderByIdResponse{
		Code:  http.StatusOK,
		Order: u,
	}
	return reply, nil
}

func (s *RechargeAdmin) UpdateOrderAckStateById(ctx context.Context, req *adminv1.UpdateOrderAckStateByIdRequest) (*adminv1.UpdateOrderAckStateByIdResponse, error) {
	if req.TransId == "" {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "transId is empty")
	}
	if req.Store == "" {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "store is empty")
	}

	store, err := pkg.StoreFromString(req.Store)
	if err != nil {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "invalid store")
	}

	if _, ok := dbv1.OrderAckState_name[req.Ack]; !ok {
		return nil, kerrors.BadRequest(xerrors.ErrAdminParamReason, "invalid ack state")
	}

	if req.Ack != int32(dbv1.OrderAckState_ORDER_ACK_STATE_GM_SUCCEEDED) && req.Ack != int32(dbv1.OrderAckState_ORDER_ACK_STATE_GM_CANCELED) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid ack state")
	}

	if err := s.uc.UpdateAckStateByID(ctx, store, req.TransId, dbv1.OrderAckState(req.Ack)); err != nil {
		return nil, errors.BadRequest("update failed", err.Error())
	}

	reply := &adminv1.UpdateOrderAckStateByIdResponse{
		Code: http.StatusOK,
	}
	return reply, nil
}

func toOrderProto(p *dbv1.OrderProto) (*adminv1.OrderProto, error) {
	p.Token = ""

	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, errors.InternalServer("json marshal failed", err.Error())
	}

	u := &adminv1.OrderProto{
		Store:       p.Store,
		TransId:     p.TransId,
		Uid:         p.Uid,
		Ack:         p.Ack,
		ProductId:   p.ProductId,
		PurchasedAt: p.PurchasedAt,
		AckAt:       p.AckAt,
		Detail:      string(bytes),
	}
	return u, nil
}
