package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	xmongo "github.com/go-pantheon/fabrica-util/data/db/mongo"
	xredis "github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/roma/app/room/internal/conf"
	"github.com/redis/go-redis/v9"
	mongo "go.mongodb.org/mongo-driver/v2/mongo"
)

type Data struct {
	Mdb *mongo.Database
	Rdb redis.UniversalClient
}

func NewData(c *conf.Data, l log.Logger) (d *Data, cleanup func(), err error) {
	var (
		mdb        *mongo.Database
		mdbCleanup func()

		rdb        redis.UniversalClient
		rdbCleanup func()
	)

	cleanup = func() {
		if mdbCleanup != nil {
			mdbCleanup()
		}

		if rdbCleanup != nil {
			rdbCleanup()
		}
	}

	mdb, mdbCleanup, err = xmongo.New(context.Background(), c.Mongo.Source, c.Mongo.Database)
	if err != nil {
		return
	}

	if c.Redis.Cluster {
		rdb, cleanup, err = xredis.NewCluster(&redis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
		if err != nil {
			return
		}
	} else {
		rdb, cleanup, err = xredis.NewStandalone(&redis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	d = &Data{
		Mdb: mdb,
		Rdb: rdb,
	}

	return d, cleanup, nil
}
