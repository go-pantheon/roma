package gate

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/system/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/system/gate/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/system/gate/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
