package security

import (
	"crypto/ed25519"
	"sync"

	"github.com/go-pantheon/fabrica-util/security/aes"
	"github.com/go-pantheon/fabrica-util/security/certificate"
	"github.com/go-pantheon/roma/mercury/internal/conf"
)

var (
	manager = &CryptoManager{}
)

type CryptoManager struct {
	mu         sync.RWMutex
	cliCertPri ed25519.PrivateKey
	svrCertPub ed25519.PublicKey
	tokenAES   *aes.Cipher
}

func Init(c *conf.Secret) (err error) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	aesCipher, err := aes.NewAESCipher([]byte(c.AccountAesKey))
	if err != nil {
		return err
	}

	cliPri, err := certificate.ImportPriFromPEM([]byte(c.ClientCertPrivateKey))
	if err != nil {
		return err
	}

	svrCert, err := certificate.ImportCertFromPEM([]byte(c.ServerCert))
	if err != nil {
		return err
	}

	if err := certificate.VerifyCert(svrCert); err != nil {
		return err
	}

	svrPub, err := certificate.ExportPubFromCert(svrCert)
	if err != nil {
		return err
	}

	manager.tokenAES = aesCipher
	manager.cliCertPri = cliPri
	manager.svrCertPub = svrPub

	return nil
}
