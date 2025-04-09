package storage

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/util/maths/i64"
	"github.com/go-pantheon/roma/pkg/util/maths/u64"
)

var _ cmds.Commandable = (*AddPackCommander)(nil)

const (
	AddPackArgItemId = "ItemID"
	AddPackArgAmount = "Amount"
)

type AddPackCommander struct {
	*cmds.BaseCommander

	storageDo *domain.StorageDomain
}

func NewAddPackCommander(uc *biz.DevUseCase, storageDo *domain.StorageDomain) *AddPackCommander {
	c := &AddPackCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Add Pack",
			"",
			[]*climsg.DevCmdArgProto{
				{
					Key: AddPackArgItemId, Def: "1",
				},
				{
					Key: AddPackArgAmount, Def: "1",
				},
			}),
		storageDo: storageDo,
	}

	uc.Register(c)
	return c
}

func (c *AddPackCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}

	var (
		packId int64
		count  uint64
	)

	if packId, err = i64.ToI64(args[AddPackArgItemId]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}
	if count, err = u64.ToU64(args[AddPackArgAmount]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	prizes, err := gamedata.TryNewPackPrize(packId, count)
	if err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	if err = c.storageDo.Add(ctx, domain.WithPacks(prizes)); err != nil {
		sc.Message = err.Error()
		return
	}

	sc.Code = climsg.SCDevExecute_Succeeded
	return
}
