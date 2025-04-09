package domain

import (
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	adv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/storage/v1"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/pkg/errors"
)

type AddOption func(o *AddOptions)

type AddOptions struct {
	Items      []*gamedata.ItemPrize
	Packs      []*gamedata.PackPrize
	Prizes     []*gamedata.Prizes
	ItemSource adv1.ItemSource
	Silent     bool
}

func WithItems(items ...*gamedata.ItemPrize) AddOption {
	return func(o *AddOptions) {
		o.Items = items
	}
}

func WithPacks(packs ...*gamedata.PackPrize) AddOption {
	return func(o *AddOptions) {
		o.Packs = packs
	}
}

func WithPrizes(prizes ...*gamedata.Prizes) AddOption {
	return func(o *AddOptions) {
		o.Prizes = prizes
	}
}

func WithItemSource(itemSource adv1.ItemSource) AddOption {
	return func(o *AddOptions) {
		o.ItemSource = itemSource
	}
}

func WithSilent(silent bool) AddOption {
	return func(o *AddOptions) {
		o.Silent = silent
	}
}

func (do *StorageDomain) Add(ctx core.Context, opts ...AddOption) error {
	if err := do.Recover(ctx); err != nil {
		do.log.WithContext(ctx).Errorf("recover failed before add. uid=%d %+v", ctx.UID(), err)
	}

	option := &AddOptions{}
	for _, opt := range opts {
		opt(option)
	}

	// check first
	if len(option.Prizes) > 0 {
		if err := do.CanAddPrizes(option.Prizes...); err != nil {
			return err
		}
	}
	if len(option.Items) > 0 {
		if err := do.CanAddItemPrizes(option.Items...); err != nil {
			return err
		}
	}
	if len(option.Packs) > 0 {
		if err := do.CanAddPackPrizes(option.Packs...); err != nil {
			return err
		}
	}

	updateInfo := object.NewUpdateInfo(ctx.Now(), object.UpdateTypeAdd)

	if len(option.Prizes) > 0 {
		if up, err := do.addPrizes(ctx, option.Prizes...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	if len(option.Items) > 0 {
		if up, err := do.addItems(ctx, option.Items...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	if len(option.Packs) > 0 {
		if up, err := do.addPacks(ctx, option.Packs...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}
	ctx.Changed()
	do.AfterUpdate(ctx, updateInfo, option.ItemSource, option.Silent)
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
		if err := storage.AddItem(data, prize.Amount()); err != nil {
			return nil, err
		}
		updateInfo.AddItem(data, prize.Amount())
	}
	ctx.Changed()
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
		if err := storage.AddPack(item.Data(), item.Amount()); err != nil {
			return nil, err
		}
		updateInfo.AddPack(item.Data(), item.Amount())
	}
	ctx.Changed()
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
	ctx.Changed()
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
