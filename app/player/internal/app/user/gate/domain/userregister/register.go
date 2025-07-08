package userregister

import (
	"maps"
	"slices"
	"sync"

	"github.com/go-pantheon/roma/pkg/data/xpg"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

var (
	register           = newUserRegister()
	allModuleKeys      = make([]life.ModuleKey, 0, 32)
	allModuleDBColumns = make(map[life.ModuleKey]string)
)

type UserRegister struct {
	newFuncs *sync.Map // map[life.ModuleKey]life.NewModuleFunc
}

func newUserRegister() *UserRegister {
	return &UserRegister{
		newFuncs: &sync.Map{},
	}
}

type Option func(key life.ModuleKey, cols map[life.ModuleKey]string)

func WithPGColumnType(columnType string) Option {
	return func(key life.ModuleKey, cols map[life.ModuleKey]string) {
		cols[key] = columnType
	}
}

func Register(key life.ModuleKey, newFunc life.NewModuleFunc, opts ...Option) {
	register.newFuncs.Store(key, newFunc)
	allModuleKeys = append(allModuleKeys, key)

	for _, opt := range opts {
		opt(key, allModuleDBColumns)
	}

	if _, ok := allModuleDBColumns[key]; !ok {
		allModuleDBColumns[key] = xpg.BytesType
	}
}

func GetPGColumnType(key life.ModuleKey) string {
	return allModuleDBColumns[key]
}

func AllModuleKeys() []life.ModuleKey {
	return slices.Clone(allModuleKeys)
}

func AllModuleDBColumns() map[life.ModuleKey]string {
	return maps.Clone(allModuleDBColumns)
}

func AllModuleDBColumnsString() map[string]string {
	cols := make(map[string]string)
	for key, colType := range allModuleDBColumns {
		cols[string(key)] = colType
	}

	return cols
}

func ForEach(fn func(key life.ModuleKey, newFunc life.NewModuleFunc)) {
	register.newFuncs.Range(func(key, value any) bool {
		fn(key.(life.ModuleKey), value.(life.NewModuleFunc))
		return true
	})
}
