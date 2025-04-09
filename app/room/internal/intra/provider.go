package intra

import (
	"github.com/go-pantheon/roma/app/room/internal/intra/filter"
	"github.com/go-pantheon/roma/app/room/internal/intra/registry"
	"github.com/go-pantheon/roma/app/room/internal/intra/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.ProviderSet,
	registry.ProviderSet,
	filter.ProviderSet,
)
