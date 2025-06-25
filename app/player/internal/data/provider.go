package data

import (
	"github.com/go-pantheon/roma/pkg/data/mongodb"
	"github.com/go-pantheon/roma/pkg/data/postgresdb"
	"github.com/go-pantheon/roma/pkg/data/redisdb"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMongoClient, NewPostgreSQLClient, NewRedisClient,
	mongodb.NewMongoDB, postgresdb.NewPostgreSQLDB, redisdb.NewRedisDB,
)
