package storage

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/go-pantheon/roma/pkg/util/maths/i64"
	"github.com/go-pantheon/roma/pkg/util/maths/u64"
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

func (c *AddItemCommander) Func(ctx core.Context, args map[string]string) (*climsg.SCDevExecute, error) {
	sc := &climsg.SCDevExecute{}

	var (
		itemId int64
		count  uint64
	)

	if itemId, err := i64.ToI64(args[AddItemArgItemId]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessagef(err, "itemId=%d", itemId)

		return sc, nil
	}
	if count, err := u64.ToU64(args[AddItemArgAmount]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessagef(err, "count=%d", count)

		return sc, nil
	}

	prizes, err := gamedata.TryNewItemPrize(itemId, count)
	if err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessage(err)

		return sc, nil
	}

	if err = c.storageDo.Add(ctx, domain.WithItems(prizes)); err != nil {
		sc.Message = life.ErrorMessage(err)
		return sc, nil
	}

	sc.Code = climsg.SCDevExecute_Succeeded

	return sc, nil
}
