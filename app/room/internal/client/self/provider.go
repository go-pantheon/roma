package self

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewSelfRouteTable)
