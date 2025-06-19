package user

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/util/maths/i64"
)

const (
	RechargeCentsArg = "Cents"
)

var _ cmds.Commandable = (*SimulateRechargeCommander)(nil)

type SimulateRechargeCommander struct {
	*cmds.BaseCommander
}

func NewSimulateRechargeCommander(uc *biz.DevUseCase) *SimulateRechargeCommander {
	cmd := &SimulateRechargeCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Simulate Recharge",
			"",
			[]*climsg.DevCmdArgProto{
				{Key: RechargeCentsArg, Def: "99"},
			}),
	}

	uc.Register(cmd)
	return cmd
}

func (c *SimulateRechargeCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}
	var (
		user  = ctx.User()
		cents int64
	)
	if cents, err = i64.ToI64(args[RechargeCentsArg]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}
	if err = user.Recharge().AddRecharge(cents); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	ctx.Changed()

	sc.Code = climsg.SCDevExecute_Succeeded
	return
}
