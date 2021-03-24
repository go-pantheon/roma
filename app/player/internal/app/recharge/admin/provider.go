package admin

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin/biz"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin/data"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin/registry"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin/service"
)

var ProviderSet = wire.NewSet(
	data.ProviderSet,
	domain.ProviderSet,
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
