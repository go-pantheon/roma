package biz

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/gamedata"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/gamedata/v1"
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
