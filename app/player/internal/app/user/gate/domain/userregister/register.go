package userregister

import (
	"sync"

	"github.com/go-pantheon/roma/pkg/universe/life"
)

var (
	register = newUserRegister()
)

type UserRegister struct {
	newFuncs *sync.Map // map[life.ModuleKey]life.NewModuleFunc
}

func newUserRegister() *UserRegister {
	return &UserRegister{
		newFuncs: &sync.Map{},
	}
}

func Register(key life.ModuleKey, newFunc life.NewModuleFunc) {
	register.newFuncs.Store(key, newFunc)
}

func ForEach(fn func(key life.ModuleKey, newFunc life.NewModuleFunc)) {
	register.newFuncs.Range(func(key, value interface{}) bool {
		fn(key.(life.ModuleKey), value.(life.NewModuleFunc))
		return true
	})
}
