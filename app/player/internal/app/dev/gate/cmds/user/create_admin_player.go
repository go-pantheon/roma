package user

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	heroobj "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain/object"
	storagedo "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	"github.com/go-pantheon/roma/pkg/universe/life"
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

func (c *CreateAdminPlayerCommander) Func(ctx core.Context, args map[string]string) (*climsg.SCDevExecute, error) {
	sc := &climsg.SCDevExecute{}

	// add all heroes
	for _, d := range gamedata.GetHeroBaseDataList() {
		r, _ := heroobj.NewHero(d.ID)
		ctx.User().HeroList().Heroes[r.Id] = r
	}

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
		sc.Message = life.ErrorMessage(err)

		return sc, nil
	}

	if err = c.storageDo.Add(ctx, storagedo.WithItems(prizes.Items()...)); err != nil {
		sc.Code = climsg.SCDevExecute_ErrUnspecified
		sc.Message = life.ErrorMessage(err)

		return sc, nil
	}

	ctx.Changed()

	err = ctx.Push(int32(climod.ModuleID_System), int32(cliseq.SystemSeq_ServerLogout), ctx.UID(),
		&climsg.SCServerLogout{
			Code: climsg.SCServerLogout_Waiting,
		},
	)
	if err != nil {
		sc.Code = climsg.SCDevExecute_ErrUnspecified
		sc.Message = life.ErrorMessage(err)

		return sc, nil
	}

	sc.Code = climsg.SCDevExecute_Succeeded

	return sc, nil
}
