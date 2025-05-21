package dbv1

import (
	"sync"
)

var (
	UserProtoPool         = newUserProtoPool()
	UserBasicProtoPool    = newUserBasicProtoPool()
	RechargeProtoPool     = newRechargeProtoPool()
	DevProtoPool          = newDevProtoPool()
	PlundersProtoPool     = newPlundersProtoPool()
	PlunderStateProtoPool = newPlunderStateProtoPool()
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
	DevProtoPool.Put(p.Dev)
	SystemProtoPool.Put(p.System)
	PlundersProtoPool.Put(p.Plunders)
	UserBasicProtoPool.Put(p.Basic)
	UserStorageProtoPool.Put(p.Storage)
	UserHeroListProtoPool.Put(p.HeroList)
	UserRoomProtoPool.Put(p.Room)

	p.Reset()
	pool.Pool.Put(p)
}

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
	RechargeProtoPool.Put(p.Recharge)

	p.Reset()
	pool.Pool.Put(p)
}

type rechargeProtoPool struct {
	sync.Pool
}

func newRechargeProtoPool() *rechargeProtoPool {
	return &rechargeProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &RechargeProto{}
			},
		},
	}
}

func (pool *rechargeProtoPool) Get() *RechargeProto {
	return pool.Pool.Get().(*RechargeProto)
}

func (pool *rechargeProtoPool) Put(p *RechargeProto) {

	p.Reset()
	pool.Pool.Put(p)
}

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

type plundersProtoPool struct {
	sync.Pool
}

func newPlundersProtoPool() *plundersProtoPool {
	return &plundersProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &PlundersProto{}
			},
		},
	}
}

func (pool *plundersProtoPool) Get() *PlundersProto {
	return pool.Pool.Get().(*PlundersProto)
}

func (pool *plundersProtoPool) Put(p *PlundersProto) {
	for _, v := range p.Plunders {
		PlunderStateProtoPool.Put(v)
	}

	p.Reset()
	pool.Pool.Put(p)
}

type plunderStateProtoPool struct {
	sync.Pool
}

func newPlunderStateProtoPool() *plunderStateProtoPool {
	return &plunderStateProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &PlunderStateProto{}
			},
		},
	}
}

func (pool *plunderStateProtoPool) Get() *PlunderStateProto {
	return pool.Pool.Get().(*PlunderStateProto)
}

func (pool *plunderStateProtoPool) Put(p *PlunderStateProto) {

	p.Reset()
	pool.Pool.Put(p)
}
