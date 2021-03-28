package data

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	servicev1 "github.com/vulcan-frame/vulcan-game/gen/api/server/gate/service/push/v1"
	"github.com/vulcan-frame/vulcan-game/pkg/universe/constants"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
)

type PushRepo struct {
	gate  servicev1.PushServiceClient
	gates []servicev1.PushServiceClient
	log   *log.Helper
}

func NewPushRepo(client servicev1.PushServiceClient, clients []servicev1.PushServiceClient, logger log.Logger) *PushRepo {
	return &PushRepo{
		gate:  client,
		gates: clients,
		log:   log.NewHelper(log.With(logger, "module", "universe/data/gate")),
	}
}

func (r *PushRepo) Push(c context.Context, uid int64, proto *servicev1.PushBody) error {
	in := &servicev1.PushRequest{
		Uid:    uid,
		Bodies: []*servicev1.PushBody{proto},
	}

	ctx, cancel := context.WithTimeout(c, constants.AsyncGRPCTimeout)
	defer cancel()

	ctx = xcontext.SetUID(ctx, uid)
	if _, err := r.gate.Push(ctx, in); err != nil {
		if !profile.IsLocal() {
			return errors.Wrapf(err, "repo push failed. oid=%d [%d-%d]", uid, proto.Mod, proto.Seq)
		}
	}
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
