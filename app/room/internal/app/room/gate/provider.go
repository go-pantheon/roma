package gate

import (
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/biz"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/data"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/registry"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
