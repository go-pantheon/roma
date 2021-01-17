package task

import (
	"github.com/go-kratos/kratos/log"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	clipkt "github.com/vulcan-frame/vulcan-game/gen/api/client/packet"
	"github.com/vulcan-frame/vulcan-game/gen/app/codec"
	"github.com/vulcan-frame/vulcan-game/mock/internal/base"
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
	GetObj(ctx *base.Context) int64
	IsExpectSC(mod climod.ModuleID, seq int32) bool
}

type Receiver interface {
	Receive(ctx *base.Context, sc *clipkt.Packet) (redirect *clipkt.Packet, done bool, err error)
}

func LogCS(l *log.Helper, cs *clipkt.Packet) {
	l.Infof("[send] %d-%d len:%d", cs.Mod, cs.Seq, len(cs.Data))
}

func LogSC(l *log.Helper, sc *clipkt.Packet) {
	if codec.IsPushSC(climod.ModuleID(sc.Mod), sc.Seq) {
		l.Infof("[recv] push %d-%d len:%d", sc.Mod, sc.Seq, len(sc.Data))
	} else {
		l.Infof("[recv] msg %d-%d len:%d", sc.Mod, sc.Seq, len(sc.Data))
	}
}
