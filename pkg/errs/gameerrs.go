package errs

import (
	"github.com/pkg/errors"
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
