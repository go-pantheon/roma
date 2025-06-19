package userregister

import (
	"github.com/go-pantheon/roma/pkg/universe/life"
)

var (
	register = newUserRegister()
)

type UserRegister struct {
	modules map[life.ModuleKey]life.Module
}

func newUserRegister() *UserRegister {
	return &UserRegister{
		modules: make(map[life.ModuleKey]life.Module),
	}
}

func Register(key life.ModuleKey, module life.Module) {
	register.modules[key] = module
}

func Get(key life.ModuleKey) life.Module {
	return register.modules[key]
}

func ForEach(fn func(key life.ModuleKey, module life.Module)) {
	for key, module := range register.modules {
		fn(key, module)
	}
}
