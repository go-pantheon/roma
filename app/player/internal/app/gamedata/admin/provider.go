package admin

import (
	"github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/gamedata/admin/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
