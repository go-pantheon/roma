package security

import (
	"testing"
	"time"

	"github.com/go-pantheon/fabrica-util/security/aes"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/gate/intra/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

//nolint:paralleltest
func TestEncryptAccountToken(t *testing.T) {
	cipher, err := aes.NewAESCipher([]byte("ihSj8Xd0CAxJoOTpdiEvts8EQNNHov3M"))
	require.NoError(t, err)

	manager.tokenAES = cipher

	token1 := &intrav1.AuthToken{
		AccountId:   1,
		Timeout:     time.Now().Add(time.Hour).Unix(),
		Unencrypted: false,
		Rand:        "random",
		Color:       "test",
		Status:      intrav1.OnlineStatus_ONLINE_STATUS_GATE,
	}

	bytes, err := proto.Marshal(token1)
	require.NoError(t, err)

	ret, err := EncryptAccountToken(bytes)
	require.NoError(t, err)

	token2 := &intrav1.AuthToken{}

	protoBytes, err := DecryptAccountToken(ret)
	require.NoError(t, err)

	require.NoError(t, proto.Unmarshal(protoBytes, token2))

	require.Equal(t, token2.AccountId, token1.AccountId)
	require.Equal(t, token2.Timeout, token1.Timeout)
	require.Equal(t, token2.Unencrypted, token1.Unencrypted)
	require.Equal(t, token2.Color, token1.Color)
	require.Equal(t, token2.Status, token1.Status)
	require.Equal(t, token2.Rand, token1.Rand)
}
