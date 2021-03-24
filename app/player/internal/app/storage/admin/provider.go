package admin

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/admin/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/admin/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/admin/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
