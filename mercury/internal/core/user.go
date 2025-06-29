package core

import (
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

type Worker interface {
	UserManager
}

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
