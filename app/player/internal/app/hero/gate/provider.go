package gate

import (
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/hero/gate/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
