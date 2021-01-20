package storage

import (
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	"github.com/vulcan-frame/vulcan-game/gamedata"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/i64"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/u64"
)

var _ cmds.Commandable = (*AddItemCommander)(nil)

const (
	AddItemArgItemId = "ItemID"
	AddItemArgAmount = "Amount"
)

type AddItemCommander struct {
	*cmds.BaseCommander

	storageDo *domain.StorageDomain
}

func NewAddItemCommander(uc *biz.DevUseCase, storageDo *domain.StorageDomain) *AddItemCommander {
	c := &AddItemCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Add Item",
			"",
			[]*climsg.DevCmdArgProto{
				{
					Key: AddItemArgItemId, Def: "1",
				},
				{
					Key: AddItemArgAmount, Def: "1",
				},
			}),
		storageDo: storageDo,
	}

	uc.Register(c)
	return c
}

func (c *AddItemCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}

	var (
		itemId int64
		count  uint64
	)

	if itemId, err = i64.ToI64(args[AddItemArgItemId]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}
	if count, err = u64.ToU64(args[AddItemArgAmount]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	prizes, err := gamedata.TryNewItemPrize(itemId, count)
	if err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	if err = c.storageDo.Add(ctx, domain.WithItems(prizes)); err != nil {
		sc.Message = err.Error()
		return
	}
	sc.Code = climsg.SCDevExecute_Succeeded
	return
}
