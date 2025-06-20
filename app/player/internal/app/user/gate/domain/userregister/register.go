package userregister

import (
	"sync"

	"github.com/go-pantheon/roma/pkg/universe/life"
)

var (
	register      = newUserRegister()
	allModuleKeys = make([]life.ModuleKey, 0, 32)
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
	allModuleKeys = append(allModuleKeys, key)
}

func AllModuleKeys() []life.ModuleKey {
	return allModuleKeys
}

func ForEach(fn func(key life.ModuleKey, newFunc life.NewModuleFunc)) {
	register.newFuncs.Range(func(key, value interface{}) bool {
		fn(key.(life.ModuleKey), value.(life.NewModuleFunc))
		return true
	})
}
