package worker

import (
	"encoding/base64"
	"time"

	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	clipkt "github.com/go-pantheon/roma/gen/api/client/packet"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/gate/intra/v1"
	"github.com/go-pantheon/roma/mercury/internal/base"
	"github.com/go-pantheon/roma/mercury/internal/base/security"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

func (w *Worker) Handshake(ctx *base.Context) error {
	cs, err := w.packHandshakeCS(ctx)
	if err != nil {
		return err
	}
	if err = w.tcpCli.Send(cs); err != nil {
		return err
	}

	timeout := time.NewTimer(5 * time.Second)
	select {
	case <-timeout.C:
		timeout.Stop()
		return errors.Errorf("worker-%d handshake timeout", w.UID())
	default:
		sc := <-w.tcpCli.Receive()
		return w.unpackHandshakeSC(ctx, sc)
	}
}

func (w *Worker) packHandshakeCS(ctx *base.Context) (result []byte, err error) {
	cs := &climsg.CSHandshake{
		Token: w.Token,
		Pub:   security.ClientPubKey,
	}

	data, err := proto.Marshal(cs)
	if err != nil {
		err = errors.Wrap(err, "CSHandshake encode failed")
		return
	}

	req := &clipkt.Packet{
		Index: ctx.SendIndex,
		Mod:   int32(climod.ModuleID_System),
		Seq:   int32(cliseq.SystemSeq_Handshake),
		Data:  data,
	}

	if result, err = proto.Marshal(req); err != nil {
		err = errors.Wrap(err, "proto marshal failed")
		return
	}
	if result, err = security.EncryptCSHandshake(result); err != nil {
		err = errors.Wrap(err, "encrypt handshake packet failed")
		return
	}

	w.log.WithContext(ctx).Infof("send handshake. len=%d token=%s", len(result), cs.Token)
	return result, nil
}

func (w *Worker) unpackHandshakeSC(ctx *base.Context, bytes []byte) (err error) {
	org, err := security.DecryptSCHandshake(bytes)
	if err != nil {
		return
	}

	reply := &clipkt.Packet{}
	if err = proto.Unmarshal(org, reply); err != nil {
		return errors.Wrapf(err, "Packet unmarshal failed")
	}

	if reply.Mod != int32(climod.ModuleID_System) || reply.Seq != int32(cliseq.SystemSeq_Handshake) {
		return errors.Errorf("not handshake msg. mod=%d seq=%d", reply.Mod, reply.Seq)
	}

	sc := &climsg.SCHandshake{}
	if err = proto.Unmarshal(reply.Data, sc); err != nil {
		return errors.Wrapf(err, "SCHandshake unmarshal failed")
	}

	if err = w.crypto.InitProtoAES(sc.Pub); err != nil {
		return
	}

	w.log.WithContext(ctx).Debugf("receive handshake. len=%d key=%s startIndex=%d", len(bytes), base64.StdEncoding.EncodeToString(sc.Pub), sc.StartIndex)

	ctx.SendIndex = sc.StartIndex
	return nil
}

func genToken(id int64) (string, error) {
	auth := &intrav1.AuthToken{
		AccountId:   id,
		Timeout:     time.Now().Add(time.Hour).Unix(),
		Unencrypted: base.Unencrypted(),
		Color:       base.Color(),
		Status:      intrav1.OnlineStatus_ONLINE_STATUS_GATE,
	}

	if base.App().StatusAdmin {
		auth.Status = intrav1.OnlineStatus_ONLINE_STATUS_ADMIN
	}

	bytes, err := proto.Marshal(auth)
	if err != nil {
		return "", errors.Wrapf(err, "AuthToken marshal failed")
	}

	token, err := security.EncryptToken(bytes)
	if err != nil {
		return "", errors.Wrapf(err, "AuthToken encrypt failed")
	}
	return token, nil
}
