package domain

import (
	"slices"

	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	adv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/storage/v1"
	"github.com/go-pantheon/roma/pkg/zerrors"
	"github.com/pkg/errors"
)

type WithAddArg func(o *addArg)

type addArg struct {
	items      []*gamedata.ItemPrize
	packs      []*gamedata.PackPrize
	prizes     []*gamedata.Prizes
	itemSource adv1.ItemSource
	silent     bool
}

func WithItems(items ...*gamedata.ItemPrize) WithAddArg {
	return func(o *addArg) {
		o.items = slices.Clone(items)
	}
}

func WithPacks(packs ...*gamedata.PackPrize) WithAddArg {
	return func(o *addArg) {
		o.packs = slices.Clone(packs)
	}
}

func WithPrizes(prizes ...*gamedata.Prizes) WithAddArg {
	return func(o *addArg) {
		o.prizes = slices.Clone(prizes)
	}
}

func WithItemSource(itemSource adv1.ItemSource) WithAddArg {
	return func(o *addArg) {
		o.itemSource = itemSource
	}
}

func WithSilent(silent bool) WithAddArg {
	return func(o *addArg) {
		o.silent = silent
	}
}

func (do *StorageDomain) Add(ctx core.Context, opts ...WithAddArg) error {
	if err := do.Recover(ctx); err != nil {
		do.log.WithContext(ctx).Errorf("recover failed before add. uid=%d %+v", ctx.UID(), err)
	}

	arg := &addArg{}
	for _, opt := range opts {
		opt(arg)
	}

	// check first
	if len(arg.prizes) > 0 {
		if err := do.CanAddPrizes(arg.prizes...); err != nil {
			return err
		}
	}

	if len(arg.items) > 0 {
		if err := do.CanAddItemPrizes(arg.items...); err != nil {
			return err
		}
	}

	if len(arg.packs) > 0 {
		if err := do.CanAddPackPrizes(arg.packs...); err != nil {
			return err
		}
	}

	updateInfo := object.NewUpdateInfo(ctx.Now(), object.UpdateTypeAdd)

	if len(arg.prizes) > 0 {
		if up, err := do.addPrizes(ctx, arg.prizes...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	if len(arg.items) > 0 {
		if up, err := do.addItems(ctx, arg.items...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	if len(arg.packs) > 0 {
		if up, err := do.addPacks(ctx, arg.packs...); err != nil {
			return err
		} else {
			_ = updateInfo.Merge(up)
		}
	}

	ctx.Changed(object.ModuleKey)
	do.AfterUpdate(ctx, updateInfo, arg.itemSource, arg.silent)

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
		return errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourceItemData")
	}

	for _, item := range items {
		if item.Data() == nil {
			return errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourceItemData")
		}

		if item.Amount() == 0 {
			return errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourceItemData id=%d, amount=%d", item.Data().Id(), item.Amount())
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
		return errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourcePackData")
	}

	for _, pack := range packs {
		if pack.Data() == nil {
			return errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourcePackData")
		}

		if pack.Amount() == 0 {
			return errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourcePackData id=%d, amount=%d", pack.Data().Id(), pack.Amount())
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
