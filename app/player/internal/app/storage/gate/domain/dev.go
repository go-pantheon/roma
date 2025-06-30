package domain

import (
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
)

func (do *StorageDomain) Clear(ctx core.Context) error {
	storage := ctx.User().Storage()
	storage.Items = make(map[int64]*object.ItemInfo)
	storage.Packs = make(map[int64]*object.PackInfo)
	storage.RecoveryInfos = make(map[int64]*object.RecoveryInfo)

	return nil
}
