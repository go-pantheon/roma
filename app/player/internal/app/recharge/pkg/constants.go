package pkg

import "github.com/pkg/errors"

type Store string

const (
	StoreGoogle = Store("GooglePlay")
	StoreApple  = Store("AppleAppStore")
)

var errStore = errors.New("store error")

func StoreFromString(s string) (Store, error) {
	switch s {
	case "GooglePlay":
		return StoreGoogle, nil
	case "AppleAppStore":
		return StoreApple, nil
	default:
		return "", errStore
	}
}
