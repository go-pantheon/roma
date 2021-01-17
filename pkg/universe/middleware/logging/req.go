package logging

import (
	"context"

	"github.com/go-kratos/kratos/log"
	jsoniter "github.com/json-iterator/go"
	"github.com/vulcan-frame/vulcan-game/gen/app/codec"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"github.com/vulcan-frame/vulcan-kit/tunnel"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
)

func Req(ctx context.Context, log *log.Helper, p tunnel.ForwardMessage, filter func(mod, seq int32) bool) {
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
