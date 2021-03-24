package service

import (
	"context"
	"net/http"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/admin/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/storage/v1"
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

	reply := &adminv1.AddItemResponse{
		Code:  http.StatusOK,
		Items: amounts,
	}
	return reply, nil
}
