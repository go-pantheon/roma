package net

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/tunnel"
)

const (
	PackLenSize = 4
	MaxBodySize = int32(1 << 14)
)

type Service interface {
	Auth(ctx context.Context, in []byte) (out []byte, ss Session, err error)
	TunnelType(mod int32) (int32, error)
	CreateTunnel(ctx context.Context, ss Session, tp int32, routerId int64, worker tunnel.Worker) (tunnel.Tunnel, error)
	OnConnected(ctx context.Context, ss Session) (err error)
	OnDisconnect(ctx context.Context, ss Session) (err error)
	Handle(ctx context.Context, ss Session, h tunnel.Holder, in []byte) (err error)
}
