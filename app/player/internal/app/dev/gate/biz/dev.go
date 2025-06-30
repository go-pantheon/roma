package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/dev/gate/cmds"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

type DevUseCase struct {
	log     *log.Helper
	mgr     *core.Manager
	cmds    map[string]map[string]cmds.Commandable
	listMsg *climsg.SCDevList
}

func NewDevUseCase(mgr *core.Manager, logger log.Logger,
) *DevUseCase {
	uc := &DevUseCase{
		log:  log.NewHelper(log.With(logger, "module", "player/dev/gate/biz")),
		mgr:  mgr,
		cmds: make(map[string]map[string]cmds.Commandable),
	}
	uc.listMsg = &climsg.SCDevList{
		Code:     climsg.SCDevList_Succeeded,
		Commands: []*climsg.DevCmdProto{},
	}
	return uc
}

func (uc *DevUseCase) Register(cmd cmds.Commandable) {
	mod := cmd.Mod()
	s := uc.cmds[mod]
	if s == nil {
		s = make(map[string]cmds.Commandable)
		uc.cmds[mod] = s
	}
	s[cmd.Cmd()] = cmd
	uc.listMsg.Commands = append(uc.listMsg.Commands, cmd.EncodeClient())
}

func (uc *DevUseCase) Execute(ctx core.Context, mod, name string, args map[string]string) (*climsg.SCDevExecute, error) {
	sc := &climsg.SCDevExecute{}

	s := uc.cmds[mod]
	if s == nil {
		sc.Code = climsg.SCDevExecute_ErrModNotExist
		sc.Message = life.Message("mod=%s", mod)

		return sc, nil
	}

	cmd := s[name]
	if cmd == nil {
		sc.Code = climsg.SCDevExecute_ErrNameNotExist
		sc.Message = life.Message("name=%s", name)

		return sc, nil
	}

	return cmd.Func(ctx, args)
}

func (uc *DevUseCase) List(ctx core.Context) (sc *climsg.SCDevList, err error) {
	return uc.listMsg, nil
}
