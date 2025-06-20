package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/trace/postgresql"
	xpg "github.com/go-pantheon/fabrica-util/data/db/postgresql"
	xredis "github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/roma/app/room/internal/conf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Data struct {
	Pdb xpg.DBPool
	Rdb redis.UniversalClient
}

func NewData(c *conf.Data, l log.Logger) (d *Data, cleanup func(), err error) {
	var (
		pdb        *pgxpool.Pool
		pdbCleanup func()

		rdb        redis.UniversalClient
		rdbCleanup func()
	)

	cleanup = func() {
		if pdbCleanup != nil {
			pdbCleanup()
		}

		if rdbCleanup != nil {
			rdbCleanup()
		}
	}

	pgConfig := xpg.NewConfig(c.Postgresql.Source, c.Postgresql.Database)
	pdb, pdbCleanup, err = postgresql.NewTracingPool(context.Background(), postgresql.DefaultPostgreSQLConfig(pgConfig))
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
		Pdb: pdb,
		Rdb: rdb,
	}

	return d, cleanup, nil
}
