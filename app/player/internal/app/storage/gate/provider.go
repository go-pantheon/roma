package gate

import (
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/biz"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/registry"
	"github.com/go-pantheon/roma/app/player/internal/app/storage/gate/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	biz.ProviderSet,
	domain.ProviderSet,
	service.ProviderSet,
	registry.ProviderSet,
)
