package dbv1

import (
	"sync"
)

var (
	SystemProtoPool = newSystemProtoPool()
	WorkerEventPool = newWorkerEventPool()
)

type systemProtoPool struct {
	sync.Pool
}

func newSystemProtoPool() *systemProtoPool {
	return &systemProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &SystemProto{}
			},
		},
	}
}

func (pool *systemProtoPool) Get() *SystemProto {
	return pool.Pool.Get().(*SystemProto)
}

func (pool *systemProtoPool) Put(p *SystemProto) {
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
