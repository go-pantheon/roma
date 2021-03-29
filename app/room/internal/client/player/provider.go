package player

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(NewUserClient, NewRouteTable, NewConn)
