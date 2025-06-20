package domain

import (
	"slices"

	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	adv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/storage/v1"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/pkg/errors"
)

type AddOpt func(o AddOption)

type AddOption struct {
	items      []*gamedata.ItemPrize
	packs      []*gamedata.PackPrize
	prizes     []*gamedata.Prizes
	itemSource adv1.ItemSource
	silent     bool
}

func WithItems(items ...*gamedata.ItemPrize) AddOpt {
	return func(o AddOption) {
		o.items = slices.Clone(items)
	}
}

func WithPacks(packs ...*gamedata.PackPrize) AddOpt {
	return func(o AddOption) {
		o.packs = slices.Clone(packs)
	}
}

func WithPrizes(prizes ...*gamedata.Prizes) AddOpt {
	return func(o AddOption) {
		o.prizes = slices.Clone(prizes)
	}
}

func WithItemSource(itemSource adv1.ItemSource) AddOpt {
	return func(o AddOption) {
		o.itemSource = itemSource
	}
}

func WithSilent(silent bool) AddOpt {
	return func(o AddOption) {
		o.silent = silent
	}
}

func (do *StorageDomain) Add(ctx core.Context, opts ...AddOpt) error {
	if err := do.Recover(ctx); err != nil {
		do.log.WithContext(ctx).Errorf("recover failed before add. uid=%d %+v", ctx.UID(), err)
	}

	option := AddOption{}
	for _, opt := range opts {
		opt(option)
	}

	// check first
	if len(option.prizes) > 0 {
		if err := do.CanAddPrizes(option.prizes...); err != nil {
			return err
		}
	}
	if len(option.items) > 0 {
		if err := do.CanAddItemPrizes(option.items...); err != nil {
			return err
		}
	}
	if len(option.packs) > 0 {
		if err := do.CanAddPackPrizes(option.packs...); err != nil {
			return err
		}
	}

	updateInfo := object.NewUpdateInfo(ctx.Now(), object.UpdateTypeAdd)

	if len(option.prizes) > 0 {
		if up, err := do.addPrizes(ctx, option.prizes...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	if len(option.items) > 0 {
		if up, err := do.addItems(ctx, option.items...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	if len(option.packs) > 0 {
		if up, err := do.addPacks(ctx, option.packs...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	ctx.Changed(object.ModuleKey)

	do.AfterUpdate(ctx, updateInfo, option.itemSource, option.silent)
	return nil
}

func (do *StorageDomain) addItems(ctx core.Context, prizes ...*gamedata.ItemPrize) (updateInfo *object.UpdateInfo, err error) {
	if err := do.CanAddItemPrizes(prizes...); err != nil {
		return nil, err
	}

	storage := ctx.User().Storage
	updateInfo = object.NewUpdateInfo(ctx.Now(), object.UpdateTypeAdd)

	for _, prize := range prizes {
		data := prize.Data()
		if err := storage().AddItem(data, prize.Amount()); err != nil {
			return nil, err
		}
		updateInfo.AddItem(data, prize.Amount())
	}
	ctx.Changed(object.ModuleKey)
	return updateInfo, nil
}

func (do *StorageDomain) CanAddItemPrizes(items ...*gamedata.ItemPrize) (err error) {
	if len(items) == 0 {
		return errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourceItemData")
	}

	for _, item := range items {
		if item.Data() == nil {
			return errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourceItemData")
		}
		if item.Amount() == 0 {
			return errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourceItemData id=%d, amount=%d", item.Data().Id(), item.Amount())
		}
	}
	return nil
}

func (do *StorageDomain) addPacks(ctx core.Context, packs ...*gamedata.PackPrize) (updateInfo *object.UpdateInfo, err error) {
	if err := do.CanAddPackPrizes(packs...); err != nil {
		return nil, err
	}

	storage := ctx.User().Storage
	updateInfo = object.NewUpdateInfo(ctx.Now(), object.UpdateTypeAdd)

	for _, item := range packs {
		if err := storage().AddPack(item.Data(), item.Amount()); err != nil {
			return nil, err
		}
		updateInfo.AddPack(item.Data(), item.Amount())
	}
	ctx.Changed(object.ModuleKey)
	return updateInfo, nil
}

func (do *StorageDomain) CanAddPackPrizes(packs ...*gamedata.PackPrize) error {
	if len(packs) == 0 {
		return errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourcePackData")
	}

	for _, pack := range packs {
		if pack.Data() == nil {
			return errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourcePackData")
		}
		if pack.Amount() == 0 {
			return errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourcePackData id=%d, amount=%d", pack.Data().Id(), pack.Amount())
		}
	}
	return nil
}

func (do *StorageDomain) addPrizes(ctx core.Context, prizes ...*gamedata.Prizes) (updateInfo *object.UpdateInfo, err error) {
	if err := do.CanAddPrizes(prizes...); err != nil {
		return nil, err
	}

	updateInfo = object.NewUpdateInfo(ctx.Now(), object.UpdateTypeAdd)
	for _, prize := range prizes {
		if up, err := do.addItems(ctx, prize.ItemPrizes()...); err != nil {
			return nil, err
		} else {
			_ = updateInfo.Merge(up)
		}
		if up, err := do.addPacks(ctx, prize.PackPrizes()...); err != nil {
			return nil, err
		} else {
			_ = updateInfo.Merge(up)
		}
	}
	ctx.Changed(object.ModuleKey)
	return updateInfo, nil
}

func (do *StorageDomain) CanAddPrizes(prizesList ...*gamedata.Prizes) error {
	for _, prizes := range prizesList {
		if err := do.CanAddItemPrizes(prizes.ItemPrizes()...); err != nil {
			return err
		}
		if err := do.CanAddPackPrizes(prizes.PackPrizes()...); err != nil {
			return err
		}
	}
	return nil
}
