package worker

import (
	"context"
	"io"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/pkg/errors"
)

func UnlogFilter(err error) bool {
	return errors.Is(err, xsync.ErrStopByTrigger) || IsConnectionError(err) || IsContextError(err) || xerrors.IsLogoutError(err) || errors.Is(err, io.EOF)
}

func IsConnectionError(err error) bool {
	return errors.Is(err, io.EOF)
}

func IsContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
