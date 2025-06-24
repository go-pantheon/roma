package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/errors"
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

func Reply(ctx context.Context, log *log.Helper, uid int64, in xnet.TunnelMessage, out []byte, delay time.Duration, filter func(mod, seq int32) bool) {
	if !profile.IsDev() {
		return
	}

	if !filter(in.GetMod(), in.GetSeq()) {
		return
	}

	var (
		tag string
		err error
	)

	p := &clipkt.Packet{
		Mod:   in.GetMod(),
		Seq:   in.GetSeq(),
		Obj:   in.GetObj(),
		Index: in.GetIndex(),
		Data:  out,
	}

	body, cserr := codec.UnmarshalSC(p)
	if cserr != nil {
		err = errors.Join(err, cserr)
	}

	str, jsonerr := jsoniter.MarshalToString(body)
	if jsonerr != nil {
		err = errors.Join(err, jsonerr)
	}

	if codec.IsPushSC(climod.ModuleID(p.Mod), p.Seq) {
		tag = "PUS"
	} else {
		tag = "REP"
	}

	msg := fmt.Sprintf("[%s] uid=%d i=%d %d-%d oid=%d delay=%d", tag, uid, p.Index, p.Mod, p.Seq, p.Obj, delay.Milliseconds())

	if err != nil {
		log.WithContext(ctx).Debugf("%s err=%s", msg, err.Error())
	} else {
		log.WithContext(ctx).Debugf("%s body=%s", msg, str)
	}
}
