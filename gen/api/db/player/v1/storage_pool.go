package dbv1

import (
	"sync"
)

var (
	UserStorageProtoPool      = newUserStorageProtoPool()
	ItemRecoveryInfoProtoPool = newItemRecoveryInfoProtoPool()
)

type userStorageProtoPool struct {
	sync.Pool
}

func newUserStorageProtoPool() *userStorageProtoPool {
	return &userStorageProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserStorageProto{}
			},
		},
	}
}

func (pool *userStorageProtoPool) Get() *UserStorageProto {
	return pool.Pool.Get().(*UserStorageProto)
}

func (pool *userStorageProtoPool) Put(p *UserStorageProto) {
	for _, v := range p.RecoveryInfos {
		ItemRecoveryInfoProtoPool.Put(v)
	}

	p.Reset()
	pool.Pool.Put(p)
}

type itemRecoveryInfoProtoPool struct {
	sync.Pool
}

func newItemRecoveryInfoProtoPool() *itemRecoveryInfoProtoPool {
	return &itemRecoveryInfoProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &ItemRecoveryInfoProto{}
			},
		},
	}
}

func (pool *itemRecoveryInfoProtoPool) Get() *ItemRecoveryInfoProto {
	return pool.Pool.Get().(*ItemRecoveryInfoProto)
}

func (pool *itemRecoveryInfoProtoPool) Put(p *ItemRecoveryInfoProto) {

	p.Reset()
	pool.Pool.Put(p)
}
