package task

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/pkg/errors"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	clipkt "github.com/vulcan-frame/vulcan-game/gen/api/client/packet"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ Taskable = (*CommonTask)(nil)

type AssertFunc func(ctx *base.Context, cs, sc proto.Message) (done bool, err error)

type CommonTask struct {
	T      Type
	Mod    climod.ModuleID
	Seq    int32
	scType reflect.Type

	CS     proto.Message
	Assert AssertFunc
}

func NewCommonTask(t Type, mod climod.ModuleID, seq int32, scType reflect.Type, cs proto.Message, assert AssertFunc) *CommonTask {
	o := &CommonTask{
		T:      t,
		Mod:    mod,
		Seq:    seq,
		scType: scType,
		CS:     cs,
		Assert: assert,
	}
	return o
}

func (t *CommonTask) IsExpectSC(mod climod.ModuleID, seq int32) bool {
	return t.Mod == mod && t.Seq == seq
}

func (t *CommonTask) Receive(ctx *base.Context, p *clipkt.Packet) (out *clipkt.Packet, done bool, err error) {
	sc, err := t.UnmarshalSC(p)
	if err != nil {
		return
	}

	slog.Info("receive message", "msg", protojson.Format(sc))
	if !t.IsExpectSC(climod.ModuleID(p.Mod), p.Seq) {
		out = p
		return
	}

	if done, err = t.Assert(ctx, t.CS, sc); err != nil {
		return
	}
	return
}

func (t *CommonTask) UnmarshalSC(p *clipkt.Packet) (sc proto.Message, err error) {
	sc = reflect.New(t.scType).Interface().(proto.Message)
	if err = proto.Unmarshal(p.Data, sc); err != nil {
		err = errors.Wrapf(err, "message unmarshal failed. %+v", t.scType)
		return
	}
	return
}

func (t *CommonTask) Type() Type {
	return t.T
}

func (t *CommonTask) CSPacket() *clipkt.Packet {
	if t.CS == nil {
		return &clipkt.Packet{
			Mod: int32(t.Mod),
			Seq: t.Seq,
		}
	}

	cs, err := proto.Marshal(t.CS)
	if err != nil {
		panic(errors.WithStack(err))
	}

	return &clipkt.Packet{
		Mod:  int32(t.Mod),
		Seq:  t.Seq,
		Data: cs,
	}
}

func (t *CommonTask) GetObj(ctx *base.Context) int64 {
	user, err := ctx.Manager.GetClientUser()
	if err != nil {
		return 0
	}

	fmt.Println(user.Basic.Name)

	switch t.Mod {
	case climod.ModuleID_Room:
		return 0
	default:
		return ctx.Manager.UID()
	}
}
