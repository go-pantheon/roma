package dbv1

import (
	"sync"
)

var (
	UserBasicProtoPool = newUserBasicProtoPool()
)

type userBasicProtoPool struct {
	sync.Pool
}

func newUserBasicProtoPool() *userBasicProtoPool {
	return &userBasicProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserBasicProto{}
			},
		},
	}
}

func (pool *userBasicProtoPool) Get() *UserBasicProto {
	return pool.Pool.Get().(*UserBasicProto)
}

func (pool *userBasicProtoPool) Put(p *UserBasicProto) {

	p.Reset()
	pool.Pool.Put(p)
}
