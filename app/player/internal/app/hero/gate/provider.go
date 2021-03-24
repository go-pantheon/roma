package gate

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/hero/gate/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
