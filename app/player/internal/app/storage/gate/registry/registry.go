package registry

import (
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
)

func NewStorageRegistrar(svc climsg.StorageServiceServer) *StorageRegistrar {
	return &StorageRegistrar{
		svc: svc,
	}
}

type StorageRegistrar struct {
	svc climsg.StorageServiceServer
}

func (r *StorageRegistrar) GrpcRegister(s *grpc.Server) {
	climsg.RegisterStorageServiceServer(s, r.svc)
	log.Infof("Register storage gate gRPC service")
}

func (r *StorageRegistrar) HttpRegister(s *http.Server) {
	climsg.RegisterStorageServiceHTTPServer(s, r.svc)
	log.Infof("Register storage gate HTTP service")
}
