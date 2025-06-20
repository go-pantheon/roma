package errs

import "github.com/go-pantheon/fabrica-util/errors"

var (
	ErrLifeWorkerStopped       = errors.Errorf("life.Worker stopped")
	ErrLifeWorkerTooManyErrors = errors.Errorf("life.Worker execute event too many errors")
)

var (
	ErrPersistNilProto     = errors.New("persist nil proto")
	ErrPersistEmptyModules = errors.New("persist empty modules")
)
