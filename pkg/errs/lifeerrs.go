package errs

import "github.com/pkg/errors"

var (
	ErrLifeWorkerStopped       = errors.Errorf("life.Worker stopped")
	ErrLifeWorkerTooManyErrors = errors.Errorf("life.Worker execute event too many errors")
)
