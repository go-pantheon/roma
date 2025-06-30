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

func LogCS(l *log.Helper, uid int64, cs *clipkt.Packet) {
	l.Infof("[REQU] uid=%d i=%d seq=<%d-%d> oid=%d len=%d", uid, cs.Index, cs.Mod, cs.Seq, cs.Obj, len(cs.Data))
}

func LogSC(l *log.Helper, uid int64, sc *clipkt.Packet) {
	if codec.IsPushSC(climod.ModuleID(sc.Mod), sc.Seq) {
		l.Infof("[PUSH] uid=%d i=%d seq=<%d-%d> oid=%d len=%d", uid, sc.Index, sc.Mod, sc.Seq, sc.Obj, len(sc.Data))
	} else {
		l.Infof("[RESP] uid=%d i=%d seq=<%d-%d> oid=%d len=%d", uid, sc.Index, sc.Mod, sc.Seq, sc.Obj, len(sc.Data))
	}
}
