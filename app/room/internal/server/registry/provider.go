package registry

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewGateRegistrars,
	NewAdminRegistrars,
	NewServiceRegistrars,
	NewServicelessUseCase,
)
