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

const (
	SubItemArgItemId = "ItemID"
	SubItemArgAmount = "Amount"
)

var _ cmds.Commandable = (*SubItemCommander)(nil)

type SubItemCommander struct {
	*cmds.BaseCommander

	storageDo *domain.StorageDomain
}

func NewSubItemCommander(uc *biz.DevUseCase, storageDo *domain.StorageDomain) *SubItemCommander {
	c := &SubItemCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Sub Item",
			"",
			[]*climsg.DevCmdArgProto{
				{
					Key: SubItemArgItemId, Def: "1",
				},
				{
					Key: SubItemArgAmount, Def: "1",
				},
			}),
		storageDo: storageDo,
	}

	uc.Register(c)
	return c
}

func (c *SubItemCommander) Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}

	var (
		itemId int64
		count  uint64
	)

	if itemId, err = i64.ToI64(args[SubItemArgItemId]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}
	if count, err = u64.ToU64(args[SubItemArgAmount]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	costs, err := gamedata.TryNewCosts(map[int64]uint64{itemId: count})
	if err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = err.Error()
		return
	}

	if err = c.storageDo.Cost(ctx, costs); err != nil {
		sc.Message = err.Error()
		return
	}

	sc.Code = climsg.SCDevExecute_Succeeded
	return
}
