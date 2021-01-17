package errs

import (
	"github.com/go-kratos/kratos/errors"
)

var (
	ErrProfileIllegal = errors.New(403, "request forbidden", "request status error")
)
