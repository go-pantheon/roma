package system

import (
	"fmt"
	"time"

	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/i64"
	"github.com/vulcan-frame/vulcan-util/xtime"
)

var _ cmds.Commandable = (*ChangeTimeCommander)(nil)

const (
	ChangeTimeArgMinutes = "Minutes"
)

type ChangeTimeCommander struct {
	*cmds.BaseCommander
}

func NewChangeTimeCommander(uc *biz.DevUseCase) *ChangeTimeCommander {
	cmd := &ChangeTimeCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Change time",
			"Only increase, 0 to reset.",
			[]*climsg.DevCmdArgProto{
				{
					Key: ChangeTimeArgMinutes, Def: "30",
				},
			}),
	}

	uc.Register(cmd)
	return cmd
}

func (c *ChangeTimeCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}

	user := ctx.User()

	var (
		minutes int64
		dur     time.Duration
	)

	if minutes, err = i64.ToI64(args[ChangeTimeArgMinutes]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	if minutes == 0 {
		user.Dev.ResetTimeOffset()
	} else {
		dur = time.Duration(minutes) * time.Minute
		user.Dev.SetTimeOffset(dur)
	}

	ctx.Changed()
	sc.Code = climsg.SCDevExecute_Succeeded
	sc.Message = fmt.Sprintf("New time is changed to %s. Offset to %.2fm", xtime.Format(ctx.Now()), user.Dev.TimeOffset().Minutes())
	return
}
