package system

import (
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-util/xtime"
)

var _ cmds.Commandable = (*ShowTimeCommander)(nil)

type ShowTimeCommander struct {
	*cmds.BaseCommander
}

func NewShowTimeCommander(uc *biz.DevUseCase) *ShowTimeCommander {
	cmd := &ShowTimeCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Show current time",
			"This time is the time of the user with the time offset",
			[]*climsg.DevCmdArgProto{}),
	}

	uc.Register(cmd)
	return cmd
}

func (c *ShowTimeCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}
	sc.Code = climsg.SCDevExecute_Succeeded
	sc.Message = xtime.Format(ctx.Now())
	return
}
