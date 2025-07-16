package broadcaster

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
	serviceName = "janus.gate.interface"
)

type Conn struct {
	*conn.Conn
}

func NewConn(logger log.Logger, rt *GateRouteTable, r registry.Discovery) (*Conn, error) {
	conn, err := conn.NewConn(serviceName, balancer.TypeReader, logger, rt, r)
	if err != nil {
		return nil, err
	}

	return &Conn{
		Conn: conn,
	}, nil
}

type GateRouteTable struct {
	routetable.ReadOnlyRouteTable
}

func NewGateRouteTable(db goredis.UniversalClient) *GateRouteTable {
	return &GateRouteTable{
		ReadOnlyRouteTable: routetable.NewReadOnlyRouteTable(redis.New(db), serviceName),
	}
}
