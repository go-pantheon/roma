package data

import (
	"github.com/go-kratos/kratos/log"
	"github.com/redis/go-redis"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	"github.com/vulcan-frame/vulcan-util/data/cache"
	"github.com/vulcan-frame/vulcan-util/data/db"
)

type Data struct {
	Rdb cache.Cacheable
}

func NewData(c *conf.Data, l log.Logger) (d *Data, cleanup func(), err error) {
	var (
		rdb        cache.Cacheable
		rdbCleanup func()
	)

	cleanup = func() {
		if rdbCleanup != nil {
			rdbCleanup()
		}
	}

	rdb, cleanup, err = cache.NewRedis(&redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		DialTimeout:  c.Redis.DialTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
	})
	
	if err != nil {
		return
	}

	d = &Data{
		Rdb: rdb,
	}
	return
}
