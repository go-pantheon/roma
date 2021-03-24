package admin

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/data"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/admin/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
