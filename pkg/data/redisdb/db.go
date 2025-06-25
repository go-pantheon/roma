package redisdb

import (
	"github.com/redis/go-redis/v9"
)

type DB struct {
	DB redis.UniversalClient
}

func NewRedisDB(rdb redis.UniversalClient) *DB {
	return &DB{DB: rdb}
}
