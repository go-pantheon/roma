package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster/domain"
	v1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
)

type BroadcasterUsecase struct {
	pubDo  *domain.PubDomain
	logger *log.Helper
}

func NewBroadcasterUsecase(pubDomain *domain.PubDomain, logger log.Logger) *BroadcasterUsecase {
	return &BroadcasterUsecase{
		pubDo:  pubDomain,
		logger: log.NewHelper(logger),
	}
}

func (uc *BroadcasterUsecase) Push(ctx context.Context, uid int64, color string, bodies []*v1.PushBody) error {
	err := uc.pubDo.Push(ctx, uid, color, bodies)
	if err != nil {
		return xerrors.APIDBFailed(err.Error())
	}
	return nil
}

func (uc *BroadcasterUsecase) Multicast(ctx context.Context, uids []int64, color string, bodies []*v1.PushBody) error {
	if err := uc.pubDo.Multicast(ctx, uids, color, bodies); err != nil {
		return xerrors.APIDBFailed(err.Error())
	}
	return nil
}

func (uc *BroadcasterUsecase) Broadcast(ctx context.Context, color string, sid int64, bodies []*v1.PushBody) error {
	if err := uc.pubDo.Broadcast(ctx, color, sid, bodies); err != nil {
		return xerrors.APIDBFailed(err.Error())
	}
	return nil
}
