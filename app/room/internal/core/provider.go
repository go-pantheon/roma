package core

import (
	"github.com/go-pantheon/roma/pkg/universe/data"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewManager, data.ProviderSet)
