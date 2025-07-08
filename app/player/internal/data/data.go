package data

import (
	"context"

	"github.com/go-pantheon/fabrica-kit/trace/tracepg"
	"github.com/go-pantheon/fabrica-util/data/db/mongo"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	redis "github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	goredis "github.com/redis/go-redis/v9"
	gomongo "go.mongodb.org/mongo-driver/v2/mongo"
)

func NewPostgreSQLClient(c *conf.Data) (pdb *pg.DB, cleanup func(), err error) {
	pgConfig := pg.NewConfig(c.Postgresql.Source, c.Postgresql.Database)
	return tracepg.NewDB(context.Background(), tracepg.DefaultPostgreSQLConfig(pgConfig))
}

func NewMongoClient(c *conf.Data) (mdb *gomongo.Database, cleanup func(), err error) {
	return mongo.New(context.Background(), c.Mongo.Source, c.Mongo.Database)
}

func NewRedisClient(c *conf.Data) (rdb goredis.UniversalClient, cleanup func(), err error) {
	if c.Redis.Cluster {
		return redis.NewCluster(&goredis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	}

	return redis.NewStandalone(&goredis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
}
