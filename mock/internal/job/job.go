package job

import (
	"github.com/vulcan-frame/vulcan-game/mock/internal/task"
)

const (
	TypeUserLogin = Type(iota)
	TypeUserReconnect
	TypeDev
	TypeStorage
	TypeSystem
)

type Type int64

type Job struct {
	T Type

	Tasks []task.Taskable
}
