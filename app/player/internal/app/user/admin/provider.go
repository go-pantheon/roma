package admin

import (
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/data"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
