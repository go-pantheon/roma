package admin

import (
	"github.com/google/wire"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/biz"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/data"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/domain"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/registry"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/app/room/admin/service"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
