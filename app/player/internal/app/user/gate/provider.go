package gate

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/data"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
