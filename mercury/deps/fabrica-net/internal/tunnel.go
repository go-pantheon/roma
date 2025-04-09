package internal

import (
	"context"
	"sync"

	"github.com/go-pantheon/fabrica-kit/tunnel"
)

type CreateTunnelFunc func(ctx context.Context, tp int32, rid int64) (tunnel.Tunnel, error)

type tunnelHolder struct {
	sync.RWMutex

	tunnelGroups map[int32]map[int64]tunnel.Tunnel // TunnelType -> oid -> tunnel
}

func newTunnelHolder() *tunnelHolder {
	th := &tunnelHolder{
		tunnelGroups: make(map[int32]map[int64]tunnel.Tunnel, 16),
	}
	return th
}

func (h *tunnelHolder) stop() {
	h.Lock()
	defer h.Unlock()

	for _, tg := range h.tunnelGroups {
		for _, t := range tg {
			t.TriggerStop()
		}
	}
}

func (h *tunnelHolder) tunnel(tp int32, oid int64) tunnel.Tunnel {
	h.RLock()
	defer h.RUnlock()

	tg, ok := h.tunnelGroups[tp]
	if !ok {
		return nil
	}

	t, ok := tg[oid]
	if !ok || t.IsStopping() {
		return nil
	}
	return t
}

func (h *tunnelHolder) createTunnel(ctx context.Context, tp int32, oid int64, create CreateTunnelFunc) (tunnel.Tunnel, error) {
	h.Lock()
	defer h.Unlock()

	tg, ok := h.tunnelGroups[tp]
	if !ok {
		tg = make(map[int64]tunnel.Tunnel, 16)
		h.tunnelGroups[tp] = tg
	}

	t, ok := tg[oid]
	if ok && !t.IsStopping() {
		return t, nil
	}

	t, err := create(ctx, tp, oid)
	if err != nil {
		return nil, err
	}

	tg[oid] = t
	return t, nil
}
