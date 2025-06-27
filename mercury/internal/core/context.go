package core

import (
	"context"

	"github.com/go-pantheon/roma/gen/api/client/message"
)

type UserManager interface {
	UserGetter
	UserSetter
	UserID
}

type UserID interface {
	UID() int64
}

type UserGetter interface {
	GetClientUser() (*climsg.UserProto, error)
}

type UserSetter interface {
	SetClientUser(p *climsg.UserProto)
}

type Context struct {
	context.Context

	SendIndex int32
	Manager   UserManager
}

func NewContext(ctx context.Context, manager UserManager) *Context {
	return &Context{
		Context: ctx,
		Manager: manager,
	}
}
