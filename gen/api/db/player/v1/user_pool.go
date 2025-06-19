package dbv1

import (
	"sync"
)

var (
	UserProtoPool = newUserProtoPool()
)

type userProtoPool struct {
	sync.Pool
}

func newUserProtoPool() *userProtoPool {
	return &userProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserProto{}
			},
		},
	}
}

func (pool *userProtoPool) Get() *UserProto {
	return pool.Pool.Get().(*UserProto)
}

func (pool *userProtoPool) Put(p *UserProto) {

	p.Reset()
	pool.Pool.Put(p)
}
