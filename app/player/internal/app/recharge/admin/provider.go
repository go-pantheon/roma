package admin

import (
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/data"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.ProviderSet,
	domain.ProviderSet,
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
