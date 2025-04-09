package gate

import (
	"github.com/go-pantheon/roma/app/player/internal/app/system/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/system/gate/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/system/gate/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
