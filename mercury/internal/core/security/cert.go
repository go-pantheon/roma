package security

import (
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/security/certificate"
)

func VerifyECDHSvrPubKey(data, sign []byte) error {
	manager.mu.RLock()
	cliPub := manager.svrCertPub
	manager.mu.RUnlock()

	if !certificate.Verify(cliPub, data, sign) {
		return errors.New("certificate Verify failed")
	}

	return nil
}

func SignECDHCliPubKey(key []byte) ([]byte, error) {
	manager.mu.RLock()
	cliCertPri := manager.cliCertPri
	manager.mu.RUnlock()

	ret, err := certificate.Sign(cliCertPri, key)
	if err != nil {
		return nil, err
	}

	return ret.Sign, nil
}
