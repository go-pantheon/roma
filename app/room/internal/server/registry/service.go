package registry

import (
	"github.com/go-kratos/kratos/transport/grpc"
	"github.com/go-kratos/kratos/transport/http"
	intra "github.com/vulcan-frame/vulcan-game/app/room/internal/intra/registry"
)

type ServiceRegistrars struct {
	Rgs []Registrar
}

func NewServiceRegistrars(
	_ *ServicelessUseCase,
	infra *intra.IntraRegistrar,
) *ServiceRegistrars {
	return &ServiceRegistrars{
		Rgs: []Registrar{infra},
	}
}

type Registrar interface {
	GrpcRegister(srv *grpc.Server)
	HttpRegister(srv *http.Server)
}
