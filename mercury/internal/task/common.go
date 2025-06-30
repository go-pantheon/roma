package task

import (
	"fmt"
	"log/slog"
	"reflect"

	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var _ Taskable = (*CommonTask)(nil)

type AssertFunc func(ctx core.Worker, cs, sc proto.Message) error

type CommonTask struct {
	T      Type
	Mod    climod.ModuleID
	Seq    int32
	scType reflect.Type

	CS     proto.Message
	Assert AssertFunc
}

func NewCommonTask(t Type, mod climod.ModuleID, seq int32, scType reflect.Type, cs proto.Message, assert AssertFunc) *CommonTask {
	return &CommonTask{
		T:      t,
		Mod:    mod,
		Seq:    seq,
		scType: scType,
		CS:     cs,
		Assert: assert,
	}
}

func (t *CommonTask) Receive(ctx core.Worker, p *clipkt.Packet) (err error) {
	sc, err := t.UnmarshalSC(p)
	if err != nil {
		return
	}

	return t.CommonAssert(ctx, climod.ModuleID(p.Mod), p.Seq, t.CS, sc)
}

func (t *CommonTask) UnmarshalSC(p *clipkt.Packet) (sc proto.Message, err error) {
	sc = reflect.New(t.scType).Interface().(proto.Message)
	if err = proto.Unmarshal(p.Data, sc); err != nil {
		slog.Error("message unmarshal failed", "error", err, "type", t.scType)
		return nil, errors.Wrapf(err, "message unmarshal failed. %+v", t.scType)
	}

	return sc, nil
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

func (t *CommonTask) GetObj(ctx core.Worker) int64 {
	user, err := ctx.GetClientUser()
	if err != nil {
		return 0
	}

	fmt.Println(user.Basic.Name)

	switch t.Mod {
	case climod.ModuleID_Room:
		return 0
	default:
		return ctx.UID()
	}
}

func (t *CommonTask) CommonAssert(ctx core.Worker, mod climod.ModuleID, seq int32, cs, sc proto.Message) error {
	if t.Mod != mod || t.Seq != seq {
		return errors.Errorf("expect <%d-%d>, but got <%d-%d>", mod, seq, t.Mod, t.Seq)
	}

	scValue := reflect.ValueOf(sc)
	getMethod := scValue.MethodByName("GetCode")

	if !getMethod.IsValid() || getMethod.Type().NumIn() != 0 || getMethod.Type().NumOut() != 1 {
		return errors.Errorf("<%d-%d> invalid sc message", mod, seq)
	}

	results := getMethod.Call(nil)
	codeValue := results[0]

	if codeValue.Kind() != reflect.Int32 {
		return errors.Errorf("<%d-%d> invalid sc message", mod, seq)
	}

	if code := codeValue.Int(); code != 1 {
		return errors.Errorf("<%d-%d> failed. code=%d", mod, seq, code)
	}

	return t.Assert(ctx, cs, sc)
}
