package errs

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyCost        = errors.New("cost is empty")
	ErrEmptyPrize       = errors.New("prize is empty")
	ErrCostInsufficient = errors.New("cost insufficient")
)
