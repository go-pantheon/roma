package gate

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConn, NewConns, NewRouteTable, NewClient, NewClients)
