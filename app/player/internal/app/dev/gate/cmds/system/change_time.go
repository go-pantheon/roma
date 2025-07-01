package system

import (
	"fmt"
	"time"

	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	devobj "github.com/go-pantheon/roma/app/player/internal/app/dev/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/go-pantheon/roma/pkg/util/maths/i64"
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

func (c *ChangeTimeCommander) Func(ctx core.Context, args map[string]string) (*climsg.SCDevExecute, error) {
	sc := &climsg.SCDevExecute{}

	user := ctx.User()

	var (
		minutes int64
		dur     time.Duration
	)

	if minutes, err := i64.ToI64(args[ChangeTimeArgMinutes]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessagef(err, "minutes=%d", minutes)

		return sc, nil
	}

	if minutes == 0 {
		user.Dev().ResetTimeOffset()
	} else {
		dur = time.Duration(minutes) * time.Minute
		user.Dev().SetTimeOffset(dur)
	}

	ctx.Changed(devobj.ModuleKey)

	sc.Code = climsg.SCDevExecute_Succeeded
	sc.Message = fmt.Sprintf("New time is changed to %s. Offset to %.2fm", xtime.Format(ctx.Now()), user.Dev().TimeOffset().Minutes())

	return sc, nil
}
