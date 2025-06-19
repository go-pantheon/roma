package dbv1

import (
	"sync"
)

var (
	PlunderListProtoPool = newPlunderListProtoPool()
	PlunderProtoPool     = newPlunderProtoPool()
)

type plunderListProtoPool struct {
	sync.Pool
}

func newPlunderListProtoPool() *plunderListProtoPool {
	return &plunderListProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &PlunderListProto{}
			},
		},
	}
}

func (pool *plunderListProtoPool) Get() *PlunderListProto {
	return pool.Pool.Get().(*PlunderListProto)
}

func (pool *plunderListProtoPool) Put(p *PlunderListProto) {
	for _, v := range p.Plunders {
		PlunderProtoPool.Put(v)
	}

	p.Reset()
	pool.Pool.Put(p)
}

type plunderProtoPool struct {
	sync.Pool
}

func newPlunderProtoPool() *plunderProtoPool {
	return &plunderProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &PlunderProto{}
			},
		},
	}
}

func (pool *plunderProtoPool) Get() *PlunderProto {
	return pool.Pool.Get().(*PlunderProto)
}

func (pool *plunderProtoPool) Put(p *PlunderProto) {

	p.Reset()
	pool.Pool.Put(p)
}
