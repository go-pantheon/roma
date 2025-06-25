package data

import (
	"github.com/go-pantheon/roma/pkg/data/mongodb"
	"github.com/go-pantheon/roma/pkg/data/redisdb"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewRedisClient, NewMongoClient,
	redisdb.NewRedisDB, mongodb.NewMongoDB,
)
