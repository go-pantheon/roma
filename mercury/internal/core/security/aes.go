package security

import (
	"encoding/base64"

	"github.com/go-pantheon/fabrica-util/errors"
)

func EncryptAccountToken(bytes []byte) (ret string, err error) {
	manager.mu.RLock()
	tokenAESCipher := manager.tokenAES
	manager.mu.RUnlock()

	secret, err := tokenAESCipher.Encrypt(bytes)
	if err != nil {
		return "", errors.Wrap(err, "aes Encrypt failed")
	}

	return base64.URLEncoding.EncodeToString(secret), nil
}

func DecryptAccountToken(str string) (ret []byte, err error) {
	manager.mu.RLock()
	tokenAESCipher := manager.tokenAES
	manager.mu.RUnlock()

	secret, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return nil, errors.Wrap(err, "base64 DecodeString failed")
	}

	ret, err = tokenAESCipher.Decrypt(secret)
	if err != nil {
		return nil, errors.Wrap(err, "aes Decrypt failed")
	}

	return ret, nil
}
