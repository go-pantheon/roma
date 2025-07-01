package userregister

import (
	"maps"
	"slices"
	"sync"

	"github.com/go-pantheon/fabrica-util/data/db/postgresql"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

var (
	register           = newUserRegister()
	allModuleKeys      = make([]life.ModuleKey, 0, 32)
	allModuleDBColumns = make(map[life.ModuleKey]postgresql.ColumnType)
)

type UserRegister struct {
	newFuncs *sync.Map // map[life.ModuleKey]life.NewModuleFunc
}

func newUserRegister() *UserRegister {
	return &UserRegister{
		newFuncs: &sync.Map{},
	}
}

type Option func(key life.ModuleKey, cols map[life.ModuleKey]postgresql.ColumnType)

func WithPGColumnType(columnType postgresql.ColumnType) Option {
	return func(key life.ModuleKey, cols map[life.ModuleKey]postgresql.ColumnType) {
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
		allModuleDBColumns[key] = postgresql.BYTEA
	}
}

func AllModuleKeys() []life.ModuleKey {
	return slices.Clone(allModuleKeys)
}

func AllModuleDBColumns() map[life.ModuleKey]postgresql.ColumnType {
	return maps.Clone(allModuleDBColumns)
}

func AllModuleDBColumnsString() map[string]string {
	cols := make(map[string]string)
	for key, colType := range allModuleDBColumns {
		cols[string(key)] = string(colType)
	}

	return cols
}

func ForEach(fn func(key life.ModuleKey, newFunc life.NewModuleFunc)) {
	register.newFuncs.Range(func(key, value any) bool {
		fn(key.(life.ModuleKey), value.(life.NewModuleFunc))
		return true
	})
}
