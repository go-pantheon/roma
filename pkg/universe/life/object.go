package life

type ModuleKey string

type NewModuleFunc func() Module

type Module interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
