package service

import (
	"context"
	"encoding/json"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xid"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/biz"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/user/v1"
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
		return nil, err
	}

	u, err := toUserProto(p)
	if err != nil {
		return nil, err
	}

	return &adminv1.GetByIdResponse{
		User: u,
	}, nil
}

func (s *UserAdmin) UserList(ctx context.Context, req *adminv1.UserListRequest) (*adminv1.UserListResponse, error) {
	cond, start, limit, err := buildGetUserListCond(req)
	if err != nil {
		return nil, err
	}

	protos, count, err := s.uc.GetList(ctx, start, limit, cond)
	if err != nil {
		return nil, err
	}

	reply := &adminv1.UserListResponse{
		Users: make([]*adminv1.UserProto, 0, len(protos)),
		Total: count,
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

func buildGetUserListCond(req *adminv1.UserListRequest) (cond *dbv1.UserProto, start, limit int64, err error) {
	start, limit = profile.PageStartLimit(req.Page, req.PageSize)

	cond = &dbv1.UserProto{}
	if req.Cond == nil {
		err = xerrors.APIParamInvalid("cond is nil")
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
		return nil, xerrors.APICodecFailed("json marshal failed").WithCause(err)
	}

	idStr, err := xid.EncodeID(p.Id)
	if err != nil {
		return nil, xerrors.APICodecFailed("id encode failed").WithCause(err)
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
