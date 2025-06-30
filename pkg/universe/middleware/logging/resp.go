package logging

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/errors"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
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

func Resp(ctx context.Context, log *log.Helper, uid int64, out xnet.TunnelMessage, delay time.Duration, filter func(mod, seq int32) bool) {
	if !profile.IsDev() {
		return
	}

	if !filter(out.GetMod(), out.GetSeq()) {
		return
	}

	var (
		tag string
		err error
	)

	body, cserr := codec.UnmarshalSC(out.GetMod(), out.GetSeq(), out.GetData())
	if cserr != nil {
		err = errors.Join(err, cserr)
	}

	str, jsonerr := jsoniter.MarshalToString(body)
	if jsonerr != nil {
		err = errors.Join(err, jsonerr)
	}

	if codec.IsPushSC(climod.ModuleID(out.GetMod()), out.GetSeq()) {
		tag = "PUSH"
	} else {
		tag = "RESP"
	}

	msg := fmt.Sprintf("[%s] uid=%d color=%s i=%d seq=<%d-%d> oid=%d delay=%.2fms", tag, uid, xcontext.Color(ctx), out.GetIndex(), out.GetMod(), out.GetSeq(), out.GetObj(), delay.Seconds()*1000)

	if err != nil {
		log.WithContext(ctx).Debugf("%s err=%s", msg, err.Error())
	} else {
		log.WithContext(ctx).Debugf("%s body=%s", msg, str)
	}
}
