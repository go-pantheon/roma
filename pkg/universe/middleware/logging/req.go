package logging

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-net/xnet"
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
	)

	uid, _ := xcontext.UID(ctx)
	oid, _ := xcontext.OID(ctx)
	body, _ := codec.UnmarshalCS(mod, seq, cs)
	str, _ := jsoniter.MarshalToString(body)

	log.WithContext(ctx).Debugf("[Req] %d-%d uid=%d oid=%d body=%s", mod, seq, uid, oid, str)
}
