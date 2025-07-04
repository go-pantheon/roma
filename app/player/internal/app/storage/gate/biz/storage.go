package biz

import (
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	plunderdo "github.com/go-pantheon/roma/app/player/internal/app/plunder/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/zerrors"
)

type StorageUseCase struct {
	log       *log.Helper
	do        *domain.StorageDomain
	plunderDo *plunderdo.PlunderDomain
}

func NewStorageUseCase(mgr *core.Manager, logger log.Logger, storageDo *domain.StorageDomain, plunderDo *plunderdo.PlunderDomain) *StorageUseCase {
	uc := &StorageUseCase{
		log:       log.NewHelper(log.With(logger, "module", "player/storage/gate/biz")),
		do:        storageDo,
		plunderDo: plunderDo,
	}

	mgr.RegisterOnCreatedEvent(uc.onCreated)
	mgr.RegisterSecondTick(uc.onSecondTick)

	return uc
}

func (uc *StorageUseCase) onCreated(ctx core.Context) error {
	prizes := gamedata.GetResourceConstantData().OnCreatedItemsItemPrizes
	if prizes != nil {
		if err := uc.do.Add(ctx, domain.WithItems(prizes.Items()...), domain.WithSilent(true)); err != nil {
			uc.log.WithContext(ctx).Errorf("add on created items failed. uid=%d %+v", ctx.UID(), err)
		}
	}

	return nil
}

func (uc *StorageUseCase) onSecondTick(ctx core.Context) error {
	return uc.do.Recover(ctx)
}

func (uc *StorageUseCase) UsePack(ctx core.Context, cs *climsg.CSUsePack) (*climsg.SCUsePack, error) {
	sc := &climsg.SCUsePack{}

	packData := gamedata.GetResourcePackData(cs.Id)
	if packData == nil {
		sc.Code = climsg.SCUsePack_ErrPackNotExist
		return sc, nil
	}

	prizes, err := uc.do.UsePack(ctx, packData)
	if err != nil {
		if errors.Is(err, zerrors.ErrStoragePackNotFound) {
			sc.Code = climsg.SCUsePack_ErrPackNotExist
		} else {
			uc.log.WithContext(ctx).Errorf("use pack failed. uid=%d %+v", ctx.UID(), err)

			sc.Code = climsg.SCUsePack_ErrUnspecified
		}

		return sc, nil
	}

	sc.Code = climsg.SCUsePack_Succeeded
	sc.Prizes = make(map[int64]int64, len(prizes.Items()))

	for _, prize := range prizes.Items() {
		sc.Prizes[prize.Data().Id()] = int64(prize.Amount())
	}

	return sc, nil
}
