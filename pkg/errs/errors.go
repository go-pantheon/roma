package errs

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-kit/xerrors"
	"github.com/vulcan-frame/vulcan-util/xsync"
)

var (
	ErrLifeWorkerStopped       = errors.Errorf("life.Worker stopped")
	ErrLifeWorkerTooManyErrors = errors.Errorf("life.Worker execute event too many errors")
)

func DontLog(err error) bool {
	return errors.Is(err, xsync.GroupStopping) || IsConnectionError(err) || IsContextError(err) || xerrors.IsLogoutError(err)
}

func IsConnectionError(err error) bool {
	return errors.Is(err, io.EOF)
}

func IsContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
