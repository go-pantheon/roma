package errs

import (
	"github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrProfileIllegal = errors.New(403, "request forbidden", "request status error")
)
