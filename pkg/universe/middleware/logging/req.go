package logging

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/gen/app/codec"
	jsoniter "github.com/json-iterator/go"
)

func Req(ctx context.Context, log *log.Helper, p xnet.TunnelMessage, filter func(mod, seq int32) bool) {
	if !profile.IsDev() {
		return
	}

	if !filter(p.GetMod(), p.GetSeq()) {
		return
	}

	var (
		mod = p.GetMod()
		seq = p.GetSeq()
		cs  = p.GetData()
		err error
	)

	body, cserr := codec.UnmarshalCS(mod, seq, cs)
	if cserr != nil {
		err = errors.Join(err, cserr)
	}

	str, jsonerr := jsoniter.MarshalToString(body)
	if jsonerr != nil {
		err = errors.Join(err, jsonerr)
	}

	msg := fmt.Sprintf("[REQU] uid=%d i=%d seq=<%d-%d> oid=%d", xcontext.UIDOrZero(ctx), p.GetIndex(), mod, seq, p.GetObj())

	if err != nil {
		log.WithContext(ctx).Debugf("%s err=%s", msg, err.Error())
	} else {
		log.WithContext(ctx).Debugf("%s body=%s", msg, str)
	}
}
