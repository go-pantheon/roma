package player

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/registry"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/data"
	"github.com/vulcan-frame/vulcan-kit/router/balancer"
	"github.com/vulcan-frame/vulcan-kit/router/conn"
	"github.com/vulcan-frame/vulcan-kit/router/routetable"
	"github.com/vulcan-frame/vulcan-kit/router/routetable/redis"
)

const (
	serviceName = "roma.player.service"
)

type Conn struct {
	*conn.Conn
}

func NewConn(logger log.Logger, rt *RouteTable, r registry.Discovery) (*Conn, error) {
	conn, err := conn.NewConn(serviceName, balancer.BalancerTypeViewer, logger, rt, r)
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
		RouteTable: routetable.NewRouteTable("player", redis.NewRedisRouteTable(d.Rdb)),
	}
}
