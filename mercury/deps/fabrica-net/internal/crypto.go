package internal

import (
	xnet "github.com/go-pantheon/fabrica-net"
	"github.com/go-pantheon/fabrica-util/security/aes"
	"github.com/pkg/errors"
)

func encrypt(ss xnet.Session, data []byte) ([]byte, error) {
	if !ss.IsCrypto() {
		return data, nil
	}

	result, err := aes.Encrypt(ss.Key(), ss.Block(), data)
	if err != nil {
		return nil, errors.Wrapf(err, "packet encrypt failed")
	}
	return result, nil
}

func decrypt(ss xnet.Session, data []byte) ([]byte, error) {
	if !ss.IsCrypto() {
		return data, nil
	}

	result, err := aes.Decrypt(ss.Key(), ss.Block(), data)
	if err != nil {
		return nil, errors.WithMessage(err, "packet decrypt failed")
	}
	return result, nil
}
