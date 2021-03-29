package gate

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/biz"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/data"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/registry"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/gate/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
