package data

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/trace/postgresql"
	xmongo "github.com/go-pantheon/fabrica-util/data/db/mongo"
	xpg "github.com/go-pantheon/fabrica-util/data/db/postgresql"
	xredis "github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewPostgreSQLClient(c *conf.Data) (pdb *pgxpool.Pool, cleanup func(), err error) {
	pgConfig := xpg.NewConfig(c.Postgresql.Source, c.Postgresql.Database)
	return postgresql.NewTracingPool(context.Background(), postgresql.DefaultPostgreSQLConfig(pgConfig))
}

func NewMongoClient(c *conf.Data) (mdb *mongo.Database, cleanup func(), err error) {
	return xmongo.New(context.Background(), c.Mongo.Source, c.Mongo.Database)
}

func NewRedisClient(c *conf.Data) (rdb redis.UniversalClient, cleanup func(), err error) {
	if c.Redis.Cluster {
		return xredis.NewCluster(&redis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	}

	return xredis.NewStandalone(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
}
