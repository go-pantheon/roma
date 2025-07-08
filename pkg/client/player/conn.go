package player

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-pantheon/fabrica-kit/router/balancer"
	"github.com/go-pantheon/fabrica-kit/router/conn"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	goredis "github.com/redis/go-redis/v9"
)

const (
	serviceName = "roma.player.service"
)

type Conn struct {
	*conn.Conn
}

func NewConn(logger log.Logger, rt *PlayerRouteTable, r registry.Discovery) (*Conn, error) {
	conn, err := conn.NewConn(serviceName, balancer.TypeReader, logger, rt, r)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Conn: conn,
	}, nil
}

type PlayerRouteTable struct {
	routetable.ReadOnlyRouteTable
}

func NewPlayerRouteTable(db goredis.UniversalClient) *PlayerRouteTable {
	return &PlayerRouteTable{
		ReadOnlyRouteTable: routetable.NewReadOnlyRouteTable(redis.New(db), serviceName),
	}
}
