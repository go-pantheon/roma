package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	servicev1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
)

type PushRepo struct {
	// TODO broadcast app client
	log *log.Helper
}

func NewPushRepo(logger log.Logger) *PushRepo {
	return &PushRepo{
		log: log.NewHelper(log.With(logger, "module", "universe/data/gate")),
	}
}

func (r *PushRepo) Push(c context.Context, uid int64, proto *servicev1.PushBody) error {
	in := &servicev1.PushRequest{
		Uid:    uid,
		Bodies: []*servicev1.PushBody{proto},
	}

	ctx, cancel := context.WithTimeout(c, constants.WorkerPushTimeout)
	defer cancel()

	ctx = xcontext.SetUID(ctx, uid)

	r.log.WithContext(ctx).Infof("mock push to %+v", in)

	return nil
}

func (r *PushRepo) Multicast(ctx context.Context, uids []int64, proto *servicev1.PushBody) error {
	r.log.WithContext(ctx).Infof("mock multicast to %d uids", len(uids))
	return nil
}

func (r *PushRepo) Broadcast(ctx context.Context, proto *servicev1.PushBody) error {
	r.log.WithContext(ctx).Infof("mock broadcast")
	return nil
}
