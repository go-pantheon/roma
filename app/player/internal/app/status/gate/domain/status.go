package domain

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain/object"
)

type StatusDomain struct {
	log *log.Helper
}

func NewStatusDomain(logger log.Logger) *StatusDomain {
	return &StatusDomain{
		log: log.NewHelper(log.With(logger, "module", "player/status/gate/domain")),
	}
}

func (d *StatusDomain) GetStatus(ctx context.Context, uid int64) (*object.Status, error) {
	return nil, nil
}
