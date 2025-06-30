package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

const MaxOffsetSec = 30

type SystemUseCase struct {
	log *log.Helper
}

func NewSystemUseCase(mgr *core.Manager, logger log.Logger) *SystemUseCase {
	uc := &SystemUseCase{
		log: log.NewHelper(log.With(logger, "module", "player/system/gate/biz")),
	}

	return uc
}

func (uc *SystemUseCase) Heartbeat(ctx core.Context, cs *climsg.CSHeartbeat) (*climsg.SCHeartbeat, error) {
	sc := &climsg.SCHeartbeat{}

	u := ctx.User()
	t := ctx.User().Now().Unix()

	if u.System().FirstHeartBeatCompleted() {
		if cs.ClientTime+MaxOffsetSec < t || cs.ClientTime-MaxOffsetSec > t {
			sc.Code = climsg.SCHeartbeat_ErrTime
			sc.ServerTime = t

			return sc, nil
		}
	}

	u.System().SetFirstHeartBeatCompleted()

	sc.ServerTime = t
	sc.Code = climsg.SCHeartbeat_Succeeded

	return sc, nil
}
