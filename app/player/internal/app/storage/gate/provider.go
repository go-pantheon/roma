package gate

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
