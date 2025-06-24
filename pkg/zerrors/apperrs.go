package zerrors

import (
	"github.com/go-pantheon/fabrica-util/errors"
	kerrors "github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrProfileIllegal = kerrors.New(403, "request forbidden", "request status error")
)
var (
	ErrAPIHandlerNotFound = errors.New("api handler not found")
)

var (
	ErrGameDataNotFound = errors.New("game data not found")
	ErrEmptyCost        = errors.New("cost is empty")
	ErrEmptyPrize       = errors.New("prize is empty")
	ErrCostInsufficient = errors.New("cost insufficient")
)

var (
	ErrStoragePackNotFound = errors.New("storage pack not found")
	ErrStorageItemNotFound = errors.New("storage item not found")
)
