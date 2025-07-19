package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster/biz"
	v1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
)

type BroadcasterService struct {
	v1.UnimplementedPushServiceServer

	uc     *biz.BroadcasterUsecase
	logger *log.Helper
}

func NewBroadcasterService(uc *biz.BroadcasterUsecase, logger log.Logger) *BroadcasterService {
	return &BroadcasterService{
		uc:     uc,
		logger: log.NewHelper(logger),
	}
}

func (s *BroadcasterService) Push(ctx context.Context, req *v1.PushRequest) (*v1.PushResponse, error) {
	err := s.uc.Push(ctx, req.Uid, req.Color, req.Bodies)
	if err != nil {
		return nil, err
	}
	return &v1.PushResponse{}, nil
}

func (s *BroadcasterService) Multicast(ctx context.Context, req *v1.MulticastRequest) (*v1.MulticastResponse, error) {
	err := s.uc.Multicast(ctx, req.Uids, req.Color, req.Bodies)
	if err != nil {
		return nil, err
	}
	return &v1.MulticastResponse{}, nil
}

func (s *BroadcasterService) Broadcast(ctx context.Context, req *v1.BroadcastRequest) (*v1.BroadcastResponse, error) {
	err := s.uc.Broadcast(ctx, req.Color, req.Sid, req.Bodies)
	if err != nil {
		return nil, err
	}
	return &v1.BroadcastResponse{}, nil
}
