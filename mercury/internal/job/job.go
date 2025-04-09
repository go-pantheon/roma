package job

import (
	"github.com/go-pantheon/roma/mercury/internal/task"
)

const (
	TypeUserLogin = Type(iota)
	TypeUserReconnect
	TypeDev
	TypeStorage
	TypeSystem
	TypeRoom
)

type Type int64

type Job struct {
	T Type

	Tasks []task.Taskable
}
