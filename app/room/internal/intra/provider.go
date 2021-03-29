package intra

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/intra/filter"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/intra/registry"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/intra/service"
)

var ProviderSet = wire.NewSet(
	service.ProviderSet,
	registry.ProviderSet,
	filter.ProviderSet,
)
