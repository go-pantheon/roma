package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
	Obj  int64  `json:"obj"`
	Mod  int32  `json:"mod"`
	Seq  int32  `json:"seq"`
}

type PubSubManager struct {
	redis  redis.UniversalClient
	logger *log.Helper
}

func NewPubSubManager(rdb redis.UniversalClient, logger log.Logger) *PubSubManager {
	return &PubSubManager{
		redis:  rdb,
		logger: log.NewHelper(logger),
	}
}

func (p *PubSubManager) Publish(ctx context.Context, topic string, data []byte) error {
	return p.redis.Publish(ctx, topic, data).Err()
}

func (p *PubSubManager) BatchPublish(ctx context.Context, topics []string, data []byte) error {
	pipe := p.redis.Pipeline()
	for _, topic := range topics {
		pipe.Publish(ctx, topic, data)
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (p *PubSubManager) BatchPublishForEach(ctx context.Context, f func(pipe redis.Pipeliner) error) error {
	pipe := p.redis.Pipeline()
	if err := f(pipe); err != nil {
		return err
	}
	_, err := pipe.Exec(ctx)
	return err
}
