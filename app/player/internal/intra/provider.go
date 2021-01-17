package intra

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/intra/filter"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/intra/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/intra/service"
)

var ProviderSet = wire.NewSet(
	service.ProviderSet,
	registry.ProviderSet,
	filter.ProviderSet,
)
