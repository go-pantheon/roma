package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/storage/v1"
)

type StorageUseCase struct {
	log *log.Helper
	do  *domain.StorageDomain
}

func NewStorageUseCase(logger log.Logger, storageDo *domain.StorageDomain) *StorageUseCase {
	uc := &StorageUseCase{
		log: log.NewHelper(log.With(logger, "module", "player/storage/admin/biz")),
		do:  storageDo,
	}

	return uc
}

func (uc *StorageUseCase) AddItems(ctx core.Context, req *adminv1.AddItemRequest) (itemAmounts map[int64]uint64, err error) {
	// TODO
	return
}
