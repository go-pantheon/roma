package logging

import (
	"context"

	"github.com/go-kratos/kratos/log"
	jsoniter "github.com/json-iterator/go"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	clipkt "github.com/vulcan-frame/vulcan-game/gen/api/client/packet"
	cliseq "github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
	"github.com/vulcan-frame/vulcan-game/gen/app/codec"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"github.com/vulcan-frame/vulcan-kit/tunnel"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
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

func Reply(ctx context.Context, log *log.Helper, reply tunnel.ForwardMessage, filter func(mod, seq int32) bool) {
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
