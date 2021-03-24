package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kratos/kratos/errors"
	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/biz"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/user/v1"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/i32"
	"github.com/vulcan-frame/vulcan-kit/xerrors"
	"github.com/vulcan-frame/vulcan-util/id"
)

const (
	pageSizeMax = 200
)

type UserAdmin struct {
	adminv1.UnimplementedUserAdminServer

	log *log.Helper
	uc  *biz.UserUseCase
}

func NewUserAdmin(logger log.Logger, uc *biz.UserUseCase) adminv1.UserAdminServer {
	return &UserAdmin{
		log: log.NewHelper(log.With(logger, "module", "player/user/admin/service")),
		uc:  uc,
	}
}

func (s *UserAdmin) GetById(ctx context.Context, req *adminv1.GetByIdRequest) (*adminv1.GetByIdResponse, error) {
	p, err := s.uc.GetByID(ctx, req.Id)
	if err != nil {
		return nil, errors.BadRequest(xerrors.ErrAdminQueryFailedReason, err.Error())
	}

	u, err := toUserProto(p)
	if err != nil {
		return nil, err
	}

	reply := &adminv1.GetByIdResponse{
		Code: http.StatusOK,
		User: u,
	}
	return reply, nil
}

func (s *UserAdmin) UserList(ctx context.Context, req *adminv1.UserListRequest) (*adminv1.UserListResponse, error) {
	cond, page, pageSize, err := buildGetUserListCond(req)
	if err != nil {
		return nil, err
	}

	protos, count, err := s.uc.GetList(ctx, i32.Max(page-1, 0)*pageSize, pageSize, cond)
	if err != nil {
		return nil, errors.BadRequest(xerrors.ErrAdminQueryFailedReason, err.Error())
	}

	reply := &adminv1.UserListResponse{
		Code:  http.StatusOK,
		Users: make([]*adminv1.UserProto, 0, len(protos)),
		Total: uint32(count),
	}

	for _, p := range protos {
		u, err := toUserProto(p)
		if err != nil {
			s.log.WithContext(ctx).Errorf("user proto convert failed. %+v", err)
			continue
		}
		reply.Users = append(reply.Users, u)
	}
	return reply, nil
}

func buildGetUserListCond(req *adminv1.UserListRequest) (cond *dbv1.UserProto, page, pageSize int32, err error) {
	if req.PageSize > pageSizeMax {
		err = errors.BadRequest(xerrors.ErrAdminParamReason, fmt.Sprintf("page size <= %d", pageSizeMax))
		return
	}

	if page = req.Page; page <= 0 {
		page = 1
	}
	if pageSize = req.PageSize; pageSize <= 0 {
		pageSize = 10
	}

	cond = &dbv1.UserProto{}
	if req.Cond == nil {
		err = errors.BadRequest(xerrors.ErrAdminParamReason, "cond is nil")
		return
	}

	if len(req.Cond.Name) > 0 {
		cond.Name = req.Cond.Name
	}
	return
}

func toUserProto(p *dbv1.UserProto) (*adminv1.UserProto, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return nil, errors.InternalServer(xerrors.ErrAdminUpdateFailedReason, err.Error())
	}

	idStr, err := id.EncodeId(p.Id)
	if err != nil {
		return nil, errors.InternalServer(xerrors.ErrAdminUpdateFailedReason, err.Error())
	}

	u := &adminv1.UserProto{
		Id:           p.Id,
		IdStr:        idStr,
		Name:         p.Name,
		CreatedAt:    p.CreatedAt,
		LoginAt:      p.LoginAt,
		LastOnlineAt: p.LastOnlineAt,
		LastOnlineIp: p.LastOnlineIp,
		Detail:       string(bytes),
	}
	return u, nil
}
