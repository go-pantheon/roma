package dbv1

import (
	"sync"
)

var (
	DevProtoPool = newDevProtoPool()
)

type devProtoPool struct {
	sync.Pool
}

func newDevProtoPool() *devProtoPool {
	return &devProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &DevProto{}
			},
		},
	}
}

func (pool *devProtoPool) Get() *DevProto {
	return pool.Pool.Get().(*DevProto)
}

func (pool *devProtoPool) Put(p *DevProto) {

	p.Reset()
	pool.Pool.Put(p)
}
