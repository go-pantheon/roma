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

const (
	SubPackArgItemId = "ItemID"
	SubPackArgAmount = "Amount"
)

var _ cmds.Commandable = (*SubPackCommander)(nil)

type SubPackCommander struct {
	*cmds.BaseCommander

	storageDo *domain.StorageDomain
}

func NewSubPackCommander(uc *biz.DevUseCase, storageDo *domain.StorageDomain) *SubPackCommander {
	c := &SubPackCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Sub Pack",
			"",
			[]*climsg.DevCmdArgProto{
				{
					Key: SubPackArgItemId, Def: "1",
				},
				{
					Key: SubPackArgAmount, Def: "1",
				},
			}),
		storageDo: storageDo,
	}

	uc.Register(c)

	return c
}

func (c *SubPackCommander) Func(ctx core.Context, args map[string]string) (*climsg.SCDevExecute, error) {
	sc := &climsg.SCDevExecute{}

	var (
		packId int64
		count  uint64
	)

	if packId, err := i64.ToI64(args[SubPackArgItemId]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessagef(err, "packId=%d", packId)

		return sc, nil
	}

	if count, err := u64.ToU64(args[SubPackArgAmount]); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessagef(err, "count=%d", count)

		return sc, nil
	}

	costs, err := gamedata.TryNewCosts(map[int64]uint64{packId: count})
	if err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessage(err)

		return sc, nil
	}

	if err = c.storageDo.Cost(ctx, costs); err != nil {
		sc.Message = life.ErrorMessage(err)
		return sc, nil
	}

	sc.Code = climsg.SCDevExecute_Succeeded

	return sc, nil
}
