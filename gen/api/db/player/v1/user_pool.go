package dbv1

import (
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	UserProtoPool       = newUserProtoPool()
	UserModuleProtoPool = newUserModuleProtoPool()
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

	for _, v := range p.Modules {
		UserModuleProtoPool.Put(v)
	}

	p.Reset()
	pool.Pool.Put(p)
}

type userModuleProtoPool struct {
	sync.Pool
}

func newUserModuleProtoPool() *userModuleProtoPool {
	return &userModuleProtoPool{
		Pool: sync.Pool{
			New: func() any {
				return &UserModuleProto{}
			},
		},
	}
}

func (pool *userModuleProtoPool) Get() *UserModuleProto {
	return pool.Pool.Get().(*UserModuleProto)
}

func (pool *userModuleProtoPool) Put(p *UserModuleProto) {
	if p.Module != nil {
		switch p.Module.(type) {
		case *UserModuleProto_Basic:
			UserBasicProtoPool.Put(p.GetBasic())
		case *UserModuleProto_Dev:
			UserDevProtoPool.Put(p.GetDev())
		case *UserModuleProto_Status:
			UserStatusProtoPool.Put(p.GetStatus())
		case *UserModuleProto_System:
			UserSystemProtoPool.Put(p.GetSystem())
		case *UserModuleProto_PlunderList:
			UserPlunderListProtoPool.Put(p.GetPlunderList())
		case *UserModuleProto_HeroList:
			UserHeroListProtoPool.Put(p.GetHeroList())
		case *UserModuleProto_Storage:
			UserStorageProtoPool.Put(p.GetStorage())
		case *UserModuleProto_Recharge:
			UserRechargeProtoPool.Put(p.GetRecharge())
		case *UserModuleProto_Room:
			UserRoomProtoPool.Put(p.GetRoom())
		default:
			log.Errorf("UserModuleProto put invalid type: %T", p.Module)
		}
	}

	p.Reset()
	pool.Pool.Put(p)
}
