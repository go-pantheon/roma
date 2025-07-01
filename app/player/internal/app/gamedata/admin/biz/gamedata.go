package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/player/admin/gamedata/v1"
)

type GamedataUseCase struct {
	log *log.Helper
}

func NewGamedataUseCase(mgr *core.Manager, logger log.Logger) *GamedataUseCase {
	uc := &GamedataUseCase{
		log: log.NewHelper(log.With(logger, "module", "player/gamedata/admin/biz")),
	}

	return uc
}

func (uc *GamedataUseCase) GetItemList(ctx context.Context) ([]*adminv1.ItemDataProto, []*adminv1.PackDataProto) {
	amounts := make([]*adminv1.ItemDataProto, 0, len(gamedata.GetResourceItemDataList()))
	packs := make([]*adminv1.PackDataProto, 0, len(gamedata.GetResourceItemDataList()))

	for _, item := range gamedata.GetResourceItemDataList() {
		amounts = append(amounts, &adminv1.ItemDataProto{Id: item.ID, Name: item.Name})
	}

	for _, pack := range gamedata.GetResourcePackDataList() {
		packs = append(packs, &adminv1.PackDataProto{Id: pack.ID, Name: pack.Name})
	}

	return amounts, packs
}
