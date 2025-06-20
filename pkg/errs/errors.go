package errs

import (
	"context"
	"io"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/fabrica-util/errors"
)

func IsUnloggableErr(err error) bool {
	return errors.Is(err, xsync.ErrStopByTrigger) || IsEOFError(err) || IsCancelError(err) || xerrors.IsLogoutError(err)
}

func IsEOFError(err error) bool {
	return errors.Is(err, io.EOF)
}

func IsCancelError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
