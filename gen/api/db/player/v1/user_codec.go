// Code generated by gen-api-db. DO NOT EDIT.

package dbv1

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func EncodeUserModuleProto(module life.Module) *UserModuleProto {
	p := module.EncodeServer()

	switch p.(type) {
	case *UserBasicProto:
		return p.(*UserBasicProto).Wrap()
	case *UserDevProto:
		return p.(*UserDevProto).Wrap()
	case *UserStatusProto:
		return p.(*UserStatusProto).Wrap()
	case *UserSystemProto:
		return p.(*UserSystemProto).Wrap()
	case *UserPlunderListProto:
		return p.(*UserPlunderListProto).Wrap()
	case *UserHeroListProto:
		return p.(*UserHeroListProto).Wrap()
	case *UserStorageProto:
		return p.(*UserStorageProto).Wrap()
	case *UserRechargeProto:
		return p.(*UserRechargeProto).Wrap()
	case *UserRoomProto:
		return p.(*UserRoomProto).Wrap()
	default:
		return nil
	}
}

func DecodeUserModuleProto(p *UserModuleProto, module life.Module) error {
	if p.Module == nil {
		return nil
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

func (x *UserBasicProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetBasic()
	mp.Module.(*UserModuleProto_Basic).Basic = x

	return mp
}

func (x *UserDevProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetDev()
	mp.Module.(*UserModuleProto_Dev).Dev = x

	return mp
}

func (x *UserStatusProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetStatus()
	mp.Module.(*UserModuleProto_Status).Status = x

	return mp
}

func (x *UserSystemProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetSystem()
	mp.Module.(*UserModuleProto_System).System = x

	return mp
}

func (x *UserPlunderListProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetPlunderList()
	mp.Module.(*UserModuleProto_PlunderList).PlunderList = x

	return mp
}

func (x *UserHeroListProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetHeroList()
	mp.Module.(*UserModuleProto_HeroList).HeroList = x

	return mp
}

func (x *UserStorageProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetStorage()
	mp.Module.(*UserModuleProto_Storage).Storage = x

	return mp
}

func (x *UserRechargeProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetRecharge()
	mp.Module.(*UserModuleProto_Recharge).Recharge = x

	return mp
}

func (x *UserRoomProto) Wrap() *UserModuleProto {
	mp := UserModuleProtoPool.GetRoom()
	mp.Module.(*UserModuleProto_Room).Room = x

	return mp
}

// UnmarshalBSON implements the bson.Unmarshaler interface for UserModuleProto.
// This is required to handle the 'oneof' field when decoding from MongoDB.
func (x *UserModuleProto) UnmarshalBSON(data []byte) error {
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return err
	}

	for key, value := range m {
		// Marshal the value back to BSON to be unmarshaled into the target struct.
		valData, err := bson.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal value for key %s", key)
		}

		switch key {
		case "basic":
			basic := UserBasicProtoPool.Get()
			if err := bson.Unmarshal(valData, basic); err != nil {
				return err
			}

			mp := userModuleProtoBasicPool.get()
			mp.Basic = basic
			x.Module = mp
		case "dev":
			dev := UserDevProtoPool.Get()
			if err := bson.Unmarshal(valData, dev); err != nil {
				return err
			}

			mp := userModuleProtoDevPool.get()
			mp.Dev = dev
			x.Module = mp
		case "status":
			status := UserStatusProtoPool.Get()
			if err := bson.Unmarshal(valData, status); err != nil {
				return err
			}

			mp := userModuleProtoStatusPool.get()
			mp.Status = status
			x.Module = mp
		case "system":
			system := UserSystemProtoPool.Get()
			if err := bson.Unmarshal(valData, system); err != nil {
				return err
			}

			mp := userModuleProtoSystemPool.get()
			mp.System = system
			x.Module = mp
		case "plunder_list":
			plunder_list := UserPlunderListProtoPool.Get()
			if err := bson.Unmarshal(valData, plunder_list); err != nil {
				return err
			}

			mp := userModuleProtoPlunderListPool.get()
			mp.PlunderList = plunder_list
			x.Module = mp
		case "hero_list":
			hero_list := UserHeroListProtoPool.Get()
			if err := bson.Unmarshal(valData, hero_list); err != nil {
				return err
			}

			mp := userModuleProtoHeroListPool.get()
			mp.HeroList = hero_list
			x.Module = mp
		case "storage":
			storage := UserStorageProtoPool.Get()
			if err := bson.Unmarshal(valData, storage); err != nil {
				return err
			}

			mp := userModuleProtoStoragePool.get()
			mp.Storage = storage
			x.Module = mp
		case "recharge":
			recharge := UserRechargeProtoPool.Get()
			if err := bson.Unmarshal(valData, recharge); err != nil {
				return err
			}

			mp := userModuleProtoRechargePool.get()
			mp.Recharge = recharge
			x.Module = mp
		case "room":
			room := UserRoomProtoPool.Get()
			if err := bson.Unmarshal(valData, room); err != nil {
				return err
			}

			mp := userModuleProtoRoomPool.get()
			mp.Room = room
			x.Module = mp
		}
		// Assuming there is only one field in the 'oneof'
		if x.Module != nil {
			break
		}
	}
	return nil
}
