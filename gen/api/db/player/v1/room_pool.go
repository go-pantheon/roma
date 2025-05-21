package dbv1

import (
	"sync"
)

var (
	UserRoomProtoPool = newUserRoomProtoPool()
)

type userRoomProtoPool struct {
	sync.Pool
}

func newUserRoomProtoPool() *userRoomProtoPool {
	return &userRoomProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserRoomProto{}
			},
		},
	}
}

func (pool *userRoomProtoPool) Get() *UserRoomProto {
	return pool.Pool.Get().(*UserRoomProto)
}

func (pool *userRoomProtoPool) Put(p *UserRoomProto) {

	p.Reset()
	pool.Pool.Put(p)
}
