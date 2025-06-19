package dbv1

import (
	"sync"
)

var (
	UserPlunderListProtoPool = newUserPlunderListProtoPool()
	UserPlunderProtoPool     = newUserPlunderProtoPool()
)

type userPlunderListProtoPool struct {
	sync.Pool
}

func newUserPlunderListProtoPool() *userPlunderListProtoPool {
	return &userPlunderListProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserPlunderListProto{}
			},
		},
	}
}

func (pool *userPlunderListProtoPool) Get() *UserPlunderListProto {
	return pool.Pool.Get().(*UserPlunderListProto)
}

func (pool *userPlunderListProtoPool) Put(p *UserPlunderListProto) {
	for _, v := range p.Plunders {
		UserPlunderProtoPool.Put(v)
	}

	p.Reset()
	pool.Pool.Put(p)
}

type userPlunderProtoPool struct {
	sync.Pool
}

func newUserPlunderProtoPool() *userPlunderProtoPool {
	return &userPlunderProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserPlunderProto{}
			},
		},
	}
}

func (pool *userPlunderProtoPool) Get() *UserPlunderProto {
	return pool.Pool.Get().(*UserPlunderProto)
}

func (pool *userPlunderProtoPool) Put(p *UserPlunderProto) {

	p.Reset()
	pool.Pool.Put(p)
}
