package room

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewConn, NewRouteTable, NewClient)
