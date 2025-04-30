package errs

import (
	"context"
	"io"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/pkg/errors"
)

func DontLog(err error) bool {
	return errors.Is(err, xsync.ErrGroupStopping) || IsConnectionError(err) || IsContextError(err) || xerrors.IsLogoutError(err)
}

func IsConnectionError(err error) bool {
	return errors.Is(err, io.EOF)
}

func IsContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
