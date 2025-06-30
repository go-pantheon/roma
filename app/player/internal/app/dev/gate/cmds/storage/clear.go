package storage

import (
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

var _ cmds.Commandable = (*ClearCommander)(nil)

type ClearCommander struct {
	*cmds.BaseCommander

	storageDo *domain.StorageDomain
}

func NewClearCommander(uc *biz.DevUseCase, storageDo *domain.StorageDomain) *ClearCommander {
	c := &ClearCommander{
		BaseCommander: cmds.NewBaseCommander(
			Mod,
			"Clear",
			"",
			[]*climsg.DevCmdArgProto{},
		),
		storageDo: storageDo,
	}

	uc.Register(c)

	return c
}

func (c *ClearCommander) Func(ctx core.Context, args map[string]string) (*climsg.SCDevExecute, error) {
	sc := &climsg.SCDevExecute{}

	if err := c.storageDo.Clear(ctx); err != nil {
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		sc.Message = life.ErrorMessage(err)

		return sc, nil
	}

	sc.Code = climsg.SCDevExecute_Succeeded

	return sc, nil
}
