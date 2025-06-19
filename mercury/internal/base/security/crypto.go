package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/security/aes"
	xrsa "github.com/go-pantheon/fabrica-util/security/rsa"
)

var (
	cipher *aes.Cipher

	serverPubKey *rsa.PublicKey

	clientPriKey *rsa.PrivateKey
	ClientPubKey []byte

	tokenAESKey []byte
)

func Init(aesKeyStr, serverPubKeyStr, clientPriKeyStr string) error {
	var (
		pubKeyBytes []byte
		err         error
	)

	cipher, err = aes.NewAESCipher(tokenAESKey)
	if err != nil {
		return err
	}

	clientPriKeyBytes, err := base64.StdEncoding.DecodeString(clientPriKeyStr)
	if err != nil {
		return errors.Wrap(err, "base64 decode client private key failed")
	}
	clientPriKey, err = x509.ParsePKCS1PrivateKey(clientPriKeyBytes)
	if err != nil {
		return errors.Wrap(err, "parse client private key failed")
	}

	ClientPubKey, err = x509.MarshalPKIXPublicKey(&clientPriKey.PublicKey)
	if err != nil {
		return errors.Wrap(err, "marshal client public key failed")
	}

	tokenAESKey = []byte(aesKeyStr)

	if pubKeyBytes, err = base64.StdEncoding.DecodeString(serverPubKeyStr); err != nil {
		return errors.Wrap(err, "base64 decode server public key failed")
	}

	var serverPubKeyInterface interface{}
	if serverPubKeyInterface, err = x509.ParsePKIXPublicKey(pubKeyBytes); err != nil {
		return errors.Wrap(err, "parse server public key failed")
	}
	serverPubKey = serverPubKeyInterface.(*rsa.PublicKey)
	return nil
}

// EncryptToken simulate account service create AuthToken
func EncryptToken(org []byte) (string, error) {
	secret, err := cipher.Encrypt(org)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(secret), nil
}

func EncryptCSHandshake(org []byte) (secret []byte, err error) {
	if secret, err = rsa.EncryptPKCS1v15(rand.Reader, serverPubKey, org); err != nil {
		err = errors.Wrap(err, "server public key encrypt failed.")
	}
	return
}

func DecryptSCHandshake(secret []byte) (org []byte, err error) {
	if org, err = xrsa.Decrypt(clientPriKey, secret); err != nil {
		err = errors.Wrapf(err, "client private key decrypt failed. %s", string(secret))
	}
	return
}

type Crypto struct {
	protoAESKey []byte
}

// InitProtoAES initialize proto's AES key and block
func (c *Crypto) InitProtoAES(key []byte) error {
	c.protoAESKey = key
	return nil
}

func (c *Crypto) EncryptProto(org []byte) (secret []byte, err error) {
	secret, err = cipher.Encrypt(org)
	if err != nil {
		return nil, err
	}
	return
}

func (c *Crypto) DecryptProto(secret []byte) (org []byte, err error) {
	org, err = cipher.Decrypt(secret)
	if err != nil {
		return nil, err
	}
	return
}
