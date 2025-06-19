package dbv1

import (
	"sync"
)

var (
	UserSystemProtoPool = newUserSystemProtoPool()
	WorkerEventPool     = newWorkerEventPool()
)

type userSystemProtoPool struct {
	sync.Pool
}

func newUserSystemProtoPool() *userSystemProtoPool {
	return &userSystemProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserSystemProto{}
			},
		},
	}
}

func (pool *userSystemProtoPool) Get() *UserSystemProto {
	return pool.Pool.Get().(*UserSystemProto)
}

func (pool *userSystemProtoPool) Put(p *UserSystemProto) {

	for _, v := range p.Events {
		WorkerEventPool.Put(v)
	}

	p.Reset()
	pool.Pool.Put(p)
}

type workerEventPool struct {
	sync.Pool
}

func newWorkerEventPool() *workerEventPool {
	return &workerEventPool{
		Pool: sync.Pool{
			New: func() any {
				return &WorkerEvent{}
			},
		},
	}
}

func (pool *workerEventPool) Get() *WorkerEvent {
	return pool.Pool.Get().(*WorkerEvent)
}

func (pool *workerEventPool) Put(p *WorkerEvent) {

	p.Reset()
	pool.Pool.Put(p)
}
