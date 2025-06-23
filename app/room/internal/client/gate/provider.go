package gate

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConn, NewConns, NewGateRouteTable, NewClient, NewClients)
