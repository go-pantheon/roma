package server

import (
	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/registry"
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/conf"
	etcdclient "go.etcd.io/etcd/client/v3"
)

func NewRegistrar(conf *conf.Registry) (registry.Registrar, error) {
	client, err := etcdclient.New(etcdclient.Config{
		Endpoints: conf.Etcd.Endpoints,
		Username:  conf.Etcd.Username,
		Password:  conf.Etcd.Password,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "create etcd client failed")
	}

	return etcd.New(client), nil
}
