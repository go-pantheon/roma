package dbv1

import (
	"sync"
)

var (
	OrderProtoPool = newOrderProtoPool()
)

type orderProtoPool struct {
	sync.Pool
}

func newOrderProtoPool() *orderProtoPool {
	return &orderProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &OrderProto{}
			},
		},
	}
}

func (pool *orderProtoPool) Get() *OrderProto {
	return pool.Pool.Get().(*OrderProto)
}

func (pool *orderProtoPool) Put(p *OrderProto) {

	p.Reset()
	pool.Pool.Put(p)
}
