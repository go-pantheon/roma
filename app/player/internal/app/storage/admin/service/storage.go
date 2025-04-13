package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/core"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/storage/v1"
)

type StorageAdmin struct {
	adminv1.UnimplementedStorageAdminServer

	mgr *core.Manager
	uc  *biz.StorageUseCase
	log *log.Helper
}

func NewStorageAdmin(logger log.Logger, mgr *core.Manager, uc *biz.StorageUseCase) adminv1.StorageAdminServer {
	return &StorageAdmin{
		log: log.NewHelper(log.With(logger, "module", "player/storage/admin/service")),
		mgr: mgr,
		uc:  uc,
	}
}

func (s *StorageAdmin) AddItem(ctx context.Context, req *adminv1.AddItemRequest) (*adminv1.AddItemResponse, error) {
	amounts, err := s.uc.AddItems(ctx.(core.Context), req)
	if err != nil {
		return nil, err
	}

	return &adminv1.AddItemResponse{
		Items: amounts,
	}, nil
}
