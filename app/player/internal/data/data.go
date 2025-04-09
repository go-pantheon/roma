package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/data/cache"
	"github.com/go-pantheon/fabrica-util/data/db"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Data struct {
	Mdb *mongo.Database
	Rdb cache.Cacheable
}

func NewData(c *conf.Data, l log.Logger) (d *Data, cleanup func(), err error) {
	var (
		mdb        *mongo.Database
		mdbCleanup func()

		rdb        cache.Cacheable
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

	mdb, mdbCleanup, err = db.NewMongo(c.Mongo.Source, c.Mongo.Database)
	if err != nil {
		return
	}

	if c.Redis.Cluster {
		rdb, cleanup, err = cache.NewRedisCluster(&redis.ClusterOptions{
			Addrs:        []string{c.Redis.Addr},
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	} else {
		rdb, cleanup, err = cache.NewRedis(&redis.Options{
			Addr:         c.Redis.Addr,
			Password:     c.Redis.Password,
			DialTimeout:  c.Redis.DialTimeout.AsDuration(),
			WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
			ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		})
	}
	if err != nil {
		return
	}

	d = &Data{
		Mdb: mdb,
		Rdb: rdb,
	}
	return
}
