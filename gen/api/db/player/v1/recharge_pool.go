package dbv1

import (
	"sync"
)

var (
	RechargeProtoPool = newRechargeProtoPool()
)

type rechargeProtoPool struct {
	sync.Pool
}

func newRechargeProtoPool() *rechargeProtoPool {
	return &rechargeProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &RechargeProto{}
			},
		},
	}
}

func (pool *rechargeProtoPool) Get() *RechargeProto {
	return pool.Pool.Get().(*RechargeProto)
}

func (pool *rechargeProtoPool) Put(p *RechargeProto) {

	p.Reset()
	pool.Pool.Put(p)
}
