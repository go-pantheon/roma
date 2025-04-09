// Code generated by gen-api. DO NOT EDIT.

package codec

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
)

func UnmarshalCS(mod, seq int32, bytes []byte) (cs proto.Message, err error) {
	switch climod.ModuleID(mod) {
	case climod.ModuleID_Dev:
		return UnmarshalCSDev(seq, bytes)
	case climod.ModuleID_Hero:
		return UnmarshalCSHero(seq, bytes)
	case climod.ModuleID_Storage:
		return UnmarshalCSStorage(seq, bytes)
	case climod.ModuleID_System:
		return UnmarshalCSSystem(seq, bytes)
	case climod.ModuleID_User:
		return UnmarshalCSUser(seq, bytes)
	case climod.ModuleID_Room:
		return UnmarshalCSRoom(seq, bytes)
	default:
		err = errors.Errorf("module not found. mod=%d", mod)
		return
	}
}

func UnmarshalSC(in *clipkt.Packet) (sc proto.Message, err error) {
    if in == nil {
		err = errors.Errorf("packet is nil")
		return
	}

	switch climod.ModuleID(in.Mod) {
	case climod.ModuleID_Dev:
		return UnmarshalSCDev(in.Seq, in.Data)
	case climod.ModuleID_Hero:
		return UnmarshalSCHero(in.Seq, in.Data)
	case climod.ModuleID_Storage:
		return UnmarshalSCStorage(in.Seq, in.Data)
	case climod.ModuleID_System:
		return UnmarshalSCSystem(in.Seq, in.Data)
	case climod.ModuleID_User:
		return UnmarshalSCUser(in.Seq, in.Data)
	case climod.ModuleID_Room:
		return UnmarshalSCRoom(in.Seq, in.Data)
	default:
		err = errors.Errorf("module not found. mod=%d", in.Mod)
		return
	}
}

func IsPushSC(mod climod.ModuleID, seq int32) bool {
	switch mod {
	case climod.ModuleID_Dev:
		return IsPushSCDev(seq)
	case climod.ModuleID_Hero:
		return IsPushSCHero(seq)
	case climod.ModuleID_Storage:
		return IsPushSCStorage(seq)
	case climod.ModuleID_System:
		return IsPushSCSystem(seq)
	case climod.ModuleID_User:
		return IsPushSCUser(seq)
	case climod.ModuleID_Room:
		return IsPushSCRoom(seq)
	default:
		return false
	}
}
