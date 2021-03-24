package service

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/gamedata/admin/biz"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/gamedata/v1"
)

type GamedataAdmin struct {
	adminv1.UnimplementedGamedataAdminServer
	sync.Once

	log *log.Helper
	uc  *biz.GamedataUseCase

	itemListMsg *adminv1.GetItemListResponse
}

func NewGamedataAdmin(logger log.Logger, uc *biz.GamedataUseCase) *GamedataAdmin {
	a := &GamedataAdmin{
		log: log.NewHelper(log.With(logger, "module", "player/gamedata/admin/service")),
		uc:  uc,
	}
	return a
}

func (s *GamedataAdmin) GetItemList(ctx context.Context, req *adminv1.GetItemListRequest) (*adminv1.GetItemListResponse, error) {
	s.Do(func() {
		items, packs := s.uc.GetItemList(ctx)
		s.itemListMsg = &adminv1.GetItemListResponse{
			Code:  http.StatusOK,
			Items: items,
			Packs: packs,
		}
	})

	return s.itemListMsg, nil
}
