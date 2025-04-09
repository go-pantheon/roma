package admin

import (
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/biz"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/data"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/domain"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/registry"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
