package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	servicev1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
	"github.com/go-pantheon/roma/pkg/client/broadcaster"
	"github.com/go-pantheon/roma/pkg/universe/constants"
)

type PushRepo struct {
	publisher *broadcaster.Publisher
	log       *log.Helper
}

func NewPushRepo(publisher *broadcaster.Publisher, logger log.Logger) *PushRepo {
	return &PushRepo{
		publisher: publisher,
		log:       log.NewHelper(log.With(logger, "module", "universe/data/push")),
	}
}

func (r *PushRepo) Push(c context.Context, uid int64, proto *servicev1.PushBody) error {
	return r.Pushs(c, uid, []*servicev1.PushBody{proto})
}

func (r *PushRepo) Pushs(c context.Context, uid int64, protos []*servicev1.PushBody) error {
	ctx, cancel := context.WithTimeout(c, constants.WorkerPushTimeout)
	defer cancel()

	err := r.publisher.Push(ctx, uid, protos)
	if err != nil {
		r.log.WithContext(ctx).Errorf("push failed for uid=%d: %v", uid, err)
		return err
	}

	r.log.WithContext(ctx).Debugf("push success for uid=%d", uid)
	return nil
}

func (r *PushRepo) Multicast(c context.Context, uids []int64, proto *servicev1.PushBody) error {
	return r.Multicasts(c, uids, []*servicev1.PushBody{proto})
}

func (r *PushRepo) Multicasts(c context.Context, uids []int64, protos []*servicev1.PushBody) error {
	ctx, cancel := context.WithTimeout(c, constants.WorkerPushTimeout)
	defer cancel()

	err := r.publisher.Multicast(ctx, uids, protos)
	if err != nil {
		r.log.WithContext(ctx).Errorf("multicast failed for uids=%v: %v", uids, err)
		return err
	}

	r.log.WithContext(ctx).Debugf("multicast success for %d uids", len(uids))
	return nil
}

func (r *PushRepo) Broadcast(c context.Context, proto *servicev1.PushBody) error {
	return r.Broadcasts(c, []*servicev1.PushBody{proto})
}

func (r *PushRepo) Broadcasts(c context.Context, protos []*servicev1.PushBody) error {
	ctx, cancel := context.WithTimeout(c, constants.WorkerPushTimeout)
	defer cancel()

	err := r.publisher.Broadcast(ctx, protos)
	if err != nil {
		r.log.WithContext(ctx).Errorf("broadcast failed: %v", err)
		return err
	}

	r.log.WithContext(ctx).Debugf("broadcast success")
	return nil
}
