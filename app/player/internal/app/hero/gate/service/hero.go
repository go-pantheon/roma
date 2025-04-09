package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

type HeroService struct {
	climsg.UnimplementedHeroServiceServer

	uc  *biz.HeroUseCase
	log *log.Helper
}

func NewHeroService(logger log.Logger, uc *biz.HeroUseCase) climsg.HeroServiceServer {
	return &HeroService{
		log: log.NewHelper(log.With(logger, "module", "player/hero/gate/service")),
		uc:  uc,
	}
}

func (s *HeroService) HeroLevelUpgrade(ctx context.Context, cs *climsg.CSHeroLevelUpgrade) (*climsg.SCHeroLevelUpgrade, error) {
	return s.uc.HeroLevelUpgrade(ctx.(core.Context), cs)
}
