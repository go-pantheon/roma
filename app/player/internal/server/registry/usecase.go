package registry

// ServicelessUseCase is used to reference UseCases that are not needed by Service, avoiding errors in wire_gen.go due to unused generated code
type ServicelessUseCase struct {
}

func NewServicelessUseCase() *ServicelessUseCase {
	return &ServicelessUseCase{}
}
