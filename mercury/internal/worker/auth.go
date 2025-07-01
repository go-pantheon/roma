package worker

import (
	"context"
	"time"

	"github.com/go-pantheon/fabrica-net/xnet"
	"github.com/go-pantheon/fabrica-util/security/ecdh"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/gate/intra/v1"
	"github.com/go-pantheon/roma/mercury/internal/core"
	"github.com/go-pantheon/roma/mercury/internal/core/security"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (w *Worker) Auth(ctx context.Context, pack xnet.Pack) (xnet.Session, error) {
	req := &clipkt.Packet{}
	if err := proto.Unmarshal(pack, req); err != nil {
		return nil, errors.Wrap(err, "unmarshal packet failed")
	}

	sc := &climsg.SCHandshake{}
	if err := proto.Unmarshal(req.Data, sc); err != nil {
		return nil, errors.Wrap(err, "unmarshal handshake failed")
	}

	if sc.Pub == nil || sc.Sign == nil {
		return nil, errors.Errorf("handshake packet is invalid")
	}

	if err := security.VerifyECDHSvrPubKey(sc.Pub, sc.Sign); err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	if sc.Timestamp < now-5 || sc.Timestamp > now+5 {
		return nil, errors.Errorf("handshake packet is expired")
	}

	svrPubBytes, err := ecdh.ParseKey(sc.Pub)
	if err != nil {
		return nil, err
	}

	var (
		ecdhInfo xnet.ECDHable
		cryptor  xnet.Cryptor
	)

	if core.Unencrypted() {
		ecdhInfo = xnet.NewUnECDH()
		cryptor = xnet.NewUnCryptor()
	} else {
		ecdhInfo, err = xnet.NewECDHInfo(w.hsInfo.cliPriv, svrPubBytes)
		if err != nil {
			return nil, err
		}

		cryptor, err = xnet.NewCryptor(ecdhInfo.SharedKey())
		if err != nil {
			return nil, err
		}
	}

	return xnet.NewSession(w.userId, core.Color(), int64(intrav1.OnlineStatus_ONLINE_STATUS_GATE),
		xnet.WithEncryptor(cryptor),
		xnet.WithECDH(ecdhInfo),
		xnet.WithStartTime(sc.Timestamp),
		xnet.WithCSIndex(sc.StartIndex-1),
	), nil
}

func (w *Worker) handshakePack(ctx context.Context, token string, cliPub []byte) (ret xnet.Pack, err error) {
	cliPubSign, err := security.SignECDHCliPubKey(cliPub)
	if err != nil {
		return nil, err
	}

	cs := &climsg.CSHandshake{
		Token:     token,
		Pub:       cliPub,
		Sign:      cliPubSign,
		ServerId:  1,
		Timestamp: time.Now().Unix(),
	}

	data, err := proto.Marshal(cs)
	if err != nil {
		return nil, errors.Wrap(err, "CSHandshake encode failed")
	}

	req := &clipkt.Packet{
		Mod:  int32(climod.ModuleID_System),
		Seq:  int32(cliseq.SystemSeq_Handshake),
		Data: data,
	}

	if ret, err = proto.Marshal(req); err != nil {
		return nil, errors.Wrap(err, "proto marshal failed")
	}

	w.log.WithContext(ctx).Infof("send handshake. len=%d token=%s", len(ret), cs.Token)

	return ret, nil
}

type HandshakeInfo struct {
	token   string
	cliPub  [32]byte
	cliPriv [32]byte
}

func newHandshakeInfo(userId int64) (*HandshakeInfo, error) {
	token, err := genToken(userId)
	if err != nil {
		return nil, err
	}

	cliPriv, cliPub, err := ecdh.GenKeyPair()
	if err != nil {
		return nil, err
	}

	return &HandshakeInfo{
		token:   token,
		cliPub:  cliPub,
		cliPriv: cliPriv,
	}, nil
}

func genToken(id int64) (string, error) {
	auth := &intrav1.AuthToken{
		AccountId:   id,
		Timeout:     time.Now().Add(time.Hour).Unix(),
		Unencrypted: core.Unencrypted(),
		Color:       core.Color(),
		Status:      intrav1.OnlineStatus_ONLINE_STATUS_GATE,
	}

	if core.AppConf().StatusAdmin {
		auth.Status = intrav1.OnlineStatus_ONLINE_STATUS_ADMIN
	}

	bytes, err := proto.Marshal(auth)
	if err != nil {
		return "", errors.Wrapf(err, "AuthToken marshal failed")
	}

	return security.EncryptAccountToken(bytes)
}
