package client

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-pantheon/roma/app/room/internal/conf"
	"github.com/pkg/errors"
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
