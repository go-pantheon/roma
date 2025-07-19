package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/broadcaster/internal/core"
	"github.com/go-pantheon/roma/app/broadcaster/internal/data"
	v1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
)

const (
	broadcastTopic = "gate.broadcast.message"
)

type PubDomain struct {
	pubsub     *data.PubSubManager
	routeTable *core.RouteTableManager
	logger     *log.Helper
}

func NewPubDomain(pubsub *data.PubSubManager, routeTable *core.RouteTableManager, logger log.Logger) *PubDomain {
	return &PubDomain{
		pubsub:     pubsub,
		routeTable: routeTable,
		logger:     log.NewHelper(logger),
	}
}

func (r *PubDomain) Push(ctx context.Context, uid int64, color string, bodies []*v1.PushBody) error {
	endpoint, err := r.routeTable.GetEndpoint(ctx, uid, color)
	if err != nil {
		return errors.Wrap(err, "failed to get endpoints")
	}

	msg := toPushMessage(uid, bodies)

	bytes, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal push message")
	}

	return r.pubsub.Publish(ctx, endpoint, bytes)
}

func (r *PubDomain) Broadcast(ctx context.Context, color string, sid int64, bodies []*v1.PushBody) error {
	msg := toBroadcastMessage(color, sid, bodies)
	bytes, err := proto.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "failed to marshal push message")
	}

	if err := r.pubsub.Publish(ctx, broadcastTopic, bytes); err != nil {
		return errors.Wrap(err, "failed to publish push message")
	}

	return nil
}

func (r *PubDomain) Multicast(ctx context.Context, uids []int64, color string, bodies []*v1.PushBody) error {
	topics, err := r.routeTable.GetEndpoints(ctx, uids, color)
	if err != nil {
		return errors.Wrap(err, "failed to get endpoints")
	}

	if err := r.pubsub.BatchPublishForEach(ctx, func(pipe redis.Pipeliner) error {
		for topic, uids := range topics {
			msg := toMulticastMessage(uids, bodies)
			bytes, err := proto.Marshal(msg)
			if err != nil {
				return errors.Wrap(err, "failed to marshal push message")
			}

			if err := pipe.Publish(ctx, topic, bytes).Err(); err != nil {
				return errors.Wrap(err, "failed to publish push message")
			}
		}
		return nil
	}); err != nil {
		return errors.Wrap(err, "failed to publish push message")
	}

	return nil
}

func toMulticastMessage(uids []int64, bodies []*v1.PushBody) (msg *v1.PushMessage) {
	return &v1.PushMessage{
		Uids:   uids,
		Bodies: bodies,
	}
}

func toBroadcastMessage(color string, sid int64, bodies []*v1.PushBody) (msg *v1.PushMessage) {
	return &v1.PushMessage{
		Color:  color,
		Sid:    sid,
		Bodies: bodies,
	}
}

func toPushMessage(uid int64, bodies []*v1.PushBody) (msg *v1.PushMessage) {
	return &v1.PushMessage{
		Uids:   []int64{uid},
		Bodies: bodies,
	}
}
