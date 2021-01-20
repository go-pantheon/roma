package cmds

import (
	"github.com/vulcan-frame/vulcan-game/app/player/internal/core"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
)

type Commandable interface {
	Mod() string
	Cmd() string
	EncodeClient() *climsg.DevCmdProto
	Func(ctx core.Context, args map[string]string) (sc *climsg.SCDevExecute, err error)
}

type BaseCommander struct {
	mod  string
	cmd  string
	desc string
	args []*climsg.DevCmdArgProto
}

func NewBaseCommander(mod string, cmd string, desc string, args []*climsg.DevCmdArgProto) *BaseCommander {
	return &BaseCommander{
		mod:  mod,
		cmd:  cmd,
		desc: desc,
		args: args,
	}
}

func (s *BaseCommander) Mod() string {
	return s.mod
}

func (s *BaseCommander) Cmd() string {
	return s.cmd
}

func (s *BaseCommander) Desc() string {
	return s.desc
}

func (s *BaseCommander) EncodeClient() *climsg.DevCmdProto {
	p := &climsg.DevCmdProto{
		Mod:  s.mod,
		Name: s.cmd,
		Desc: s.desc,
		Args: s.args,
	}
	return p
}
