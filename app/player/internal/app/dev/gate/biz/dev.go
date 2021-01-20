package biz

import (
	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/dev/gate/cmds"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
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

func (uc *DevUseCase) Execute(ctx core.Context, mod, name string, args map[string]string) (sc *climsg.SCDevExecute, err error) {
	sc = &climsg.SCDevExecute{}

	s := uc.cmds[mod]
	if s == nil {
		err = errors.Errorf("dev mod not exist. mod=%s", mod)
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		return
	}
	cmd := s[name]
	if cmd == nil {
		err = errors.WithMessagef(err, "dev name not exist. name=%s", name)
		sc.Code = climsg.SCDevExecute_ErrArgFormat
		return
	}

	return cmd.Func(ctx, args)
}

func (uc *DevUseCase) List(ctx core.Context) (sc *climsg.SCDevList, err error) {
	return uc.listMsg, nil
}
