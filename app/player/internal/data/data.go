package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/trace/postgresql"
	xmongo "github.com/go-pantheon/fabrica-util/data/db/mongo"
	xpg "github.com/go-pantheon/fabrica-util/data/db/postgresql"
	xredis "github.com/go-pantheon/fabrica-util/data/redis"
	"github.com/go-pantheon/roma/app/player/internal/conf"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Data struct {
	Rdb redis.UniversalClient
	Pdb xpg.DBPool
	Mdb *mongo.Database
}

func NewData(c *conf.Data, l log.Logger) (d *Data, cleanup func(), err error) {
	var (
		rdb        redis.UniversalClient
		rdbCleanup func()

		pdb        *pgxpool.Pool
		pdbCleanup func()

		mdb        *mongo.Database
		mdbCleanup func()
	)

	cleanup = func() {
		if rdbCleanup != nil {
			rdbCleanup()
		}
		if pdbCleanup != nil {
			pdbCleanup()
		}
		if mdbCleanup != nil {
			mdbCleanup()
		}
	}

	mdb, mdbCleanup, err = xmongo.New(context.Background(), c.Mongo.Source, c.Mongo.Database)
	if err != nil {
		return
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
	} else {
		rdb, cleanup, err = xredis.NewStandalone(&redis.Options{
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
		Rdb: rdb,
		Pdb: pdb,
		Mdb: mdb,
	}

	return d, cleanup, nil
}
