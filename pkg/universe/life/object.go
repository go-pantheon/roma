package life

type ModuleKey string

type Module interface {
	Register()
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}
