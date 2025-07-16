package broadcaster

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConn, NewGateRouteTable, NewClient)
