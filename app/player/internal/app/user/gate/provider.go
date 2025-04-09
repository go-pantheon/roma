package gate

import (
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/data"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
	data.ProviderSet,
)
