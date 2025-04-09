package registry

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	intra "github.com/go-pantheon/roma/app/player/internal/intra/registry"
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
