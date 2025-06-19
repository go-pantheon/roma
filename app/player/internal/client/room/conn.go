package room

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-pantheon/fabrica-kit/router/balancer"
	"github.com/go-pantheon/fabrica-kit/router/conn"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/router/routetable/redis"
	"github.com/go-pantheon/roma/app/player/internal/data"
)

const (
	serviceName = "roma.room.service"
)

type Conn struct {
	*conn.Conn
}

func NewConn(logger log.Logger, rt *RouteTable, r registry.Discovery) (*Conn, error) {
	conn, err := conn.NewConn(serviceName, balancer.TypeReader, logger, rt, r)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Conn: conn,
	}, nil
}

type RouteTable struct {
	routetable.RouteTable
}

func NewRouteTable(d *data.Data) *RouteTable {
	return &RouteTable{
		RouteTable: routetable.NewRouteTable("room", redis.New(d.Rdb)),
	}
}
