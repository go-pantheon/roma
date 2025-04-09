package admin

import (
	"github.com/go-pantheon/roma/app/player/internal/app/storage/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/admin/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/admin/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
