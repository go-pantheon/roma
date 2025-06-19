package dbv1

import (
	"sync"
)

var (
	UserDevProtoPool = newUserDevProtoPool()
)

type userDevProtoPool struct {
	sync.Pool
}

func newUserDevProtoPool() *userDevProtoPool {
	return &userDevProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserDevProto{}
			},
		},
	}
}

func (pool *userDevProtoPool) Get() *UserDevProto {
	return pool.Pool.Get().(*UserDevProto)
}

func (pool *userDevProtoPool) Put(p *UserDevProto) {

	p.Reset()
	pool.Pool.Put(p)
}
