package security

import (
	"crypto/ed25519"
	"testing"

	"github.com/go-pantheon/fabrica-util/security/certificate"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestSignECDHCliPubKey(t *testing.T) {
	cliPub := initTest(t)

	data := []byte("test")

	sign, err := SignECDHCliPubKey(data)
	require.NoError(t, err)

	verified := certificate.Verify(cliPub, data, sign)
	require.True(t, verified)
}

func initTest(t *testing.T) (cliPub []byte) {
	pub, pri, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	cliPub = pub
	manager.cliCertPri = pri

	return cliPub
}
