package broadcaster

import (
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster/biz"
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster/domain"
	"github.com/go-pantheon/roma/app/broadcaster/internal/app/broadcaster/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.ProviderSet,
	biz.ProviderSet,
	domain.ProviderSet,
)
