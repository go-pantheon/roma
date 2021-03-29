package client

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd"
	"github.com/go-kratos/kratos/registry"
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/room/internal/conf"
	etcdclient "go.etcd.io/etcd/client/v3"
)

func NewDiscovery(conf *conf.Registry) (registry.Discovery, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: conf.Etcd.Endpoints,
		Username:  conf.Etcd.Username,
		Password:  conf.Etcd.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "new etcdclient failed")
	}

	r := etcd.New(client)
	return r, nil
}
