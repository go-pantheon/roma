package dbv1

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

func EncodeUserModuleProto(module life.Module) (*UserModuleProto, error) {
	p := module.EncodeServer()
	mp := UserModuleProtoPool.Get()

	switch p.(type) {
	case *UserBasicProto:
		mp.Module = &UserModuleProto_Basic{
			Basic: p.(*UserBasicProto),
		}
		return mp, nil
	case *UserDevProto:
		mp.Module = &UserModuleProto_Dev{
			Dev: p.(*UserDevProto),
		}
		return mp, nil
	case *UserStatusProto:
		mp.Module = &UserModuleProto_Status{
			Status: p.(*UserStatusProto),
		}
		return mp, nil
	case *UserSystemProto:
		mp.Module = &UserModuleProto_System{
			System: p.(*UserSystemProto),
		}
		return mp, nil
	case *UserPlunderListProto:
		mp.Module = &UserModuleProto_PlunderList{
			PlunderList: p.(*UserPlunderListProto),
		}
		return mp, nil
	case *UserHeroListProto:
		mp.Module = &UserModuleProto_HeroList{
			HeroList: p.(*UserHeroListProto),
		}
		return mp, nil
	case *UserStorageProto:
		mp.Module = &UserModuleProto_Storage{
			Storage: p.(*UserStorageProto),
		}
		return mp, nil
	case *UserRechargeProto:
		mp.Module = &UserModuleProto_Recharge{
			Recharge: p.(*UserRechargeProto),
		}
		return mp, nil
	case *UserRoomProto:
		mp.Module = &UserModuleProto_Room{
			Room: p.(*UserRoomProto),
		}
		return mp, nil
	default:
		return nil, errors.Errorf("UserModuleProto encode invalid type: %T", module)
	}
}

func DecodeUserModuleProto(p *UserModuleProto, module life.Module) error {
	if p.Module == nil {
		return errors.New("UserModuleProto.Module is nil")
	}

	switch p.Module.(type) {
	case *UserModuleProto_Basic:
		return module.DecodeServer(p.GetBasic())
	case *UserModuleProto_Dev:
		return module.DecodeServer(p.GetDev())
	case *UserModuleProto_Status:
		return module.DecodeServer(p.GetStatus())
	case *UserModuleProto_System:
		return module.DecodeServer(p.GetSystem())
	case *UserModuleProto_PlunderList:
		return module.DecodeServer(p.GetPlunderList())
	case *UserModuleProto_HeroList:
		return module.DecodeServer(p.GetHeroList())
	case *UserModuleProto_Storage:
		return module.DecodeServer(p.GetStorage())
	case *UserModuleProto_Recharge:
		return module.DecodeServer(p.GetRecharge())
	case *UserModuleProto_Room:
		return module.DecodeServer(p.GetRoom())
	default:
		return errors.Errorf("UserModuleProto decode invalid type: %T", p.Module)
	}
}
