package user

import (
	"github.com/go-kratos/kratos/log"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds"
	storagedo "github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/gamedata"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	cliseq "github.com/vulcan-frame/vulcan-game/gen/api/client/sequence"
)

var _ cmds.Commandable = (*CreateAdminPlayerCommander)(nil)

type CreateAdminPlayerCommander struct {
	*cmds.BaseCommander

	storageDo *storagedo.StorageDomain
}

func NewAdminPlayerCommander(uc *biz.DevUseCase, storageDo *storagedo.StorageDomain) *CreateAdminPlayerCommander {
	cmd := &CreateAdminPlayerCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Generate Admin Player",
			"",
			[]*climsg.DevCmdArgProto{}),
		storageDo: storageDo,
	}

	uc.Register(cmd)
	return cmd
}

func (c *CreateAdminPlayerCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}

	// add all necessary items
	amounts := make(map[int64]uint64)
	for _, i := range []int64{1, 2, 3} {
		if d := gamedata.GetResourceItemData(i); d != nil {
			amounts[d.ID] = 100000000
		}
	}
	prizes, err := gamedata.TryNewItemPrizes(amounts)
	if err != nil {
		log.Errorf("failed to create admin player. %+v", err)
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	if err = c.storageDo.Add(ctx, storagedo.WithItems(prizes.Items()...)); err != nil {
		log.Errorf("failed to create admin player. %+v", err)
		sc.Code = climsg.SCDevExecute_ErrUnspecified
		sc.Message = err.Error()
		return
	}

	ctx.Changed()

	_ = ctx.Reply(climod.ModuleID_System, int32(cliseq.SystemSeq_ServerLogout), ctx.UID(), &climsg.SCServerLogout{
		Code: climsg.SCServerLogout_Waiting,
	})

	sc.Code = climsg.SCDevExecute_Succeeded
	return
}
