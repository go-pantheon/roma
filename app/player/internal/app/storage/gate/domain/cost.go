package domain

import (
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/domain/object"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/gamedata"
	"github.com/vulcan-frame/vulcan-game/pkg/errs"
)

func (do *StorageDomain) Cost(ctx core.Context, toCosts *gamedata.Costs) (err error) {
	if err = do.Recover(ctx); err != nil {
		do.log.WithContext(ctx).Errorf("recover failed before cost. uid=%d %+v", ctx.UID(), err)
	}

	if err = do.CanCost(ctx, toCosts); err != nil {
		return
	}

	storage := ctx.User().Storage
	updateInfo := object.NewUpdateInfo(ctx.Now(), object.UpdateTypeSub)

	toCosts.Walk(func(itemCost *gamedata.ItemCost) bool {
		item := storage.Items[itemCost.Data().Id()]
		if item == nil {
			err = errors.Wrapf(errs.ErrCostInsufficient, "Data=ResourceItemData id=%d", itemCost.Data().Id())
			return false
		}
		item.Sub(itemCost.Amount())
		updateInfo.AddItem(itemCost.Data(), itemCost.Amount())
		return true
	})
	if err != nil {
		return err
	}

	ctx.Changed()
	do.AfterUpdate(ctx, updateInfo, 0, false)
	return
}

func (do *StorageDomain) CanCost(ctx core.Context, toCosts *gamedata.Costs) (err error) {
	storage := ctx.User().Storage

	toCosts.Walk(func(itemCost *gamedata.ItemCost) bool {
		item, ok := storage.Items[itemCost.Data().Id()]
		if !ok {
			err = errors.Wrapf(errs.ErrCostInsufficient, "Data=ResourceItemData id=%d", itemCost.Data().Id())
			return false
		}
		if item.Amount() < itemCost.Amount() {
			err = errors.Wrapf(errs.ErrCostInsufficient, "Data=ResourceItemData id=%d", itemCost.Data().Id())
			return false
		}
		return true
	})

	if err != nil {
		return err
	}
	return nil
}
