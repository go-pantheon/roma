package registry

import (
	"github.com/go-kratos/kratos/log"
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	adminv1 "github.com/vulcan-frame/vulcan-game/gen/api/server/player/admin/storage/v1"
)

func NewStorageRegistrar(svc adminv1.StorageAdminServer) *StorageRegistrar {
	return &StorageRegistrar{
		svc: svc,
	}
}

type StorageRegistrar struct {
	svc adminv1.StorageAdminServer
}

func (r *StorageRegistrar) GrpcRegister(s *grpc.Server) {
	adminv1.RegisterStorageAdminServer(s, r.svc)
	log.Infof("Register storage admin gRPC service")
}

func (r *StorageRegistrar) HttpRegister(s *http.Server) {
	adminv1.RegisterStorageAdminHTTPServer(s, r.svc)
	log.Infof("Register storage admin HTTP service")
}
