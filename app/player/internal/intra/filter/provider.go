package filter

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewHttpFilter, NewGrpcFilter)
