package domain

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	adv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/storage/v1"
	"github.com/go-pantheon/roma/pkg/zerrors"
)

func NewStorageDomain(logger log.Logger) *StorageDomain {
	return &StorageDomain{
		log: log.NewHelper(log.With(logger, "module", "player/storage/gate/domain")),
	}
}

type StorageDomain struct {
	log *log.Helper
}

func (do *StorageDomain) AfterUpdate(ctx core.Context, updateInfo *object.UpdateInfo, src adv1.ItemSource, silent bool) {
	if updateInfo == nil {
		return
	}

	var (
		itemIds   = make([]int64, 0, 16)
		packIds   = make([]int64, 0, 16)
		amountMsg = &climsg.SCPushItemUpdated{
			Items: make(map[int64]uint64, 16),
			Packs: make(map[int64]uint64, 16),
		}
	)

	storage := ctx.User().Storage()

	updateInfo.WalkItem(func(d *gamedata.ResourceItemData, amount uint64) {
		id := d.Id()
		itemIds = append(itemIds, id)
		if item := storage.Items[id]; item != nil {
			amountMsg.Items[id] = item.Amount()
		} else {
			amountMsg.Items[id] = 0
		}
	})

	updateInfo.WalkPack(func(d *gamedata.ResourcePackData, amount uint64) {
		id := d.Id()
		packIds = append(packIds, id)
		if pack := storage.Packs[id]; pack != nil {
			amountMsg.Packs[id] = pack.Amount()
		} else {
			amountMsg.Packs[id] = 0
		}
	})

	_ = ctx.ProductPreparedEvent(core.WorkerEventTypeStorageItemUpdated, itemIds...)
	_ = ctx.ProductPreparedEvent(core.WorkerEventTypeStoragePackUpdated, packIds...)

	ctx.Changed(object.ModuleKey)

	if !silent {
		_ = ctx.Reply(climod.ModuleID_Storage, int32(cliseq.StorageSeq_PushItemUpdated), ctx.UID(), amountMsg)
	}
}

func (do *StorageDomain) UsePack(ctx core.Context, packData *gamedata.ResourcePackData) (*gamedata.ItemPrizes, error) {
	storage := ctx.User().Storage()

	pack := storage.Packs[packData.Id()]
	if pack == nil {
		return nil, zerrors.ErrStoragePackNotFound
	}
	if pack.Amount() == 0 {
		return nil, zerrors.ErrStoragePackNotFound
	}

	if err := storage.SubPack(packData, pack.Amount()); err != nil {
		return nil, err
	}

	prizes := gamedata.NewEmptyItemPrizes()
	for i := 0; i < int(pack.Amount()); i++ {
		prize := packData.Rand()
		if gamedata.IsItemPrizesValid(prize) {
			prizes = prizes.CloneWith(prize)
		}
	}

	if len(prizes.Items()) == 0 {
		do.log.WithContext(ctx).Errorf("pack prize is empty. uid=%d pack=%d", ctx.UID(), packData.Id())
		return prizes, nil
	}

	if err := do.Add(ctx, WithItems(prizes.Items()...)); err != nil {
		return nil, err
	}

	return prizes, nil
}

func (do *StorageDomain) Recover(ctx core.Context) error {
	storage := ctx.User().Storage()
	ctime := ctx.Now()

	prizeList := make([]*gamedata.ItemPrize, 0, len(storage.RecoveryInfos))
	for _, rec := range storage.RecoveryInfos {
		toAdd := rec.Recover(ctime)
		if toAdd > 0 {
			prize, err := gamedata.TryNewItemPrize(rec.Id, toAdd)
			if err != nil {
				return err
			}
			prizeList = append(prizeList, prize)
		}
	}

	if len(prizeList) == 0 {
		return nil
	}

	if err := do.Add(ctx, WithItems(prizeList...), WithSilent(true)); err != nil {
		return err
	}
	return nil
}
