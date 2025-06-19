package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-net/xnet"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	"github.com/go-pantheon/roma/gen/app/codec"
	jsoniter "github.com/json-iterator/go"
)

func EmptyFilter(mod, seq int32) bool {
	return true
}

func DefaultFilter(mod, seq int32) bool {
	switch mod {
	case int32(climod.ModuleID_System):
		switch seq {
		case int32(cliseq.SystemSeq_Heartbeat):
			return false
		}
	}
	return true
}

func Reply(ctx context.Context, log *log.Helper, reply xnet.TunnelMessage, filter func(mod, seq int32) bool) {
	if !profile.IsDev() {
		return
	}

	if !filter(reply.GetMod(), reply.GetSeq()) {
		return
	}

	uid, _ := xcontext.UID(ctx)
	oid, _ := xcontext.OID(ctx)

	p := &clipkt.Packet{
		Mod:  reply.GetMod(),
		Seq:  reply.GetSeq(),
		Obj:  reply.GetObj(),
		Data: reply.GetData(),
	}

	body, _ := codec.UnmarshalSC(p)
	str, _ := jsoniter.MarshalToString(body)

	var tag string
	if codec.IsPushSC(climod.ModuleID(p.Mod), p.Seq) {
		tag = "Push"
	} else {
		tag = "Reply"
	}
	log.WithContext(ctx).Debugf("[%s] %d-%d uid=%d oid=%d body=%s", tag, p.Mod, p.Seq, uid, oid, str)
}
