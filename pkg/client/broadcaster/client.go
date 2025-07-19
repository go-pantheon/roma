package broadcaster

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/xcontext"
	servicev1 "github.com/go-pantheon/roma/gen/api/server/broadcaster/service/push/v1"
	"github.com/pkg/errors"
)

type Publisher struct {
	conn servicev1.PushServiceClient
}

func NewPublisher(conn *Conn) *Publisher {
	return &Publisher{
		conn: servicev1.NewPushServiceClient(conn.ClientConnInterface),
	}
}

func (c *Publisher) Push(ctx context.Context, uid int64, bodies []*servicev1.PushBody) error {
	req := &servicev1.PushRequest{
		Uid:    uid,
		Bodies: bodies,
	}

	_, err := c.conn.Push(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "push message failed for uid=%d", uid)
	}

	return nil
}

func (c *Publisher) Multicast(ctx context.Context, uids []int64, bodies []*servicev1.PushBody) error {
	req := &servicev1.MulticastRequest{
		Color:  xcontext.Color(ctx),
		Uids:   uids,
		Bodies: bodies,
	}

	_, err := c.conn.Multicast(ctx, req)
	if err != nil {
		return errors.Wrapf(err, "multicast message failed for uids=%v", uids)
	}

	return nil
}

func (c *Publisher) Broadcast(ctx context.Context, bodies []*servicev1.PushBody) error {
	req := &servicev1.BroadcastRequest{
		Color:  xcontext.Color(ctx),
		Sid:    xcontext.SIDOrZero(ctx),
		Bodies: bodies,
	}

	_, err := c.conn.Broadcast(ctx, req)
	if err != nil {
		return errors.Wrap(err, "broadcast message failed")
	}

	return nil
}
