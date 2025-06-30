package task

import (
	"github.com/go-kratos/kratos/v2/log"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	"github.com/go-pantheon/roma/gen/app/codec"
	"github.com/go-pantheon/roma/mercury/internal/core"
)

const (
	TypeCommon = Type(iota)
	TypeLogin
	TypeHeartbeat
)

type Type int64

type Taskable interface {
	Receiver

	Type() Type
	CSPacket() *clipkt.Packet
	GetObj(ctx core.Worker) int64
}

type Receiver interface {
	Receive(ctx core.Worker, sc *clipkt.Packet) (err error)
}

func LogCS(l *log.Helper, uid int64, p *clipkt.Packet) {
	if cs, err := codec.UnmarshalCS(p.Mod, p.Seq, p.Data); err != nil {
		l.Errorf("[REQU] uid=%d i=%d seq=<%d-%d> oid=%d err=%s", uid, p.Index, p.Mod, p.Seq, p.Obj, err)
	} else {
		l.Infof("[REQU] uid=%d i=%d seq=<%d-%d> oid=%d body=%+v", uid, p.Index, p.Mod, p.Seq, p.Obj, cs)
	}
}

func LogSC(l *log.Helper, uid int64, p *clipkt.Packet) {
	var tag string

	if codec.IsPushSC(climod.ModuleID(p.Mod), p.Seq) {
		tag = "PUSH"
	} else {
		tag = "RESP"
	}

	if sc, err := codec.UnmarshalSC(p.Mod, p.Seq, p.Data); err != nil {
		l.Errorf("[%s] uid=%d i=%d seq=<%d-%d> oid=%d err=%s", tag, uid, p.Index, p.Mod, p.Seq, p.Obj, err)
	} else {
		l.Infof("[%s] uid=%d i=%d seq=<%d-%d> oid=%d body=%+v", tag, uid, p.Index, p.Mod, p.Seq, p.Obj, sc)
	}
}
