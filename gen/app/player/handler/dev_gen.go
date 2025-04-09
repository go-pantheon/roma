// Code generated by gen-api. DO NOT EDIT.

package handler

import (
	"context"
	"google.golang.org/protobuf/proto"
	"github.com/go-pantheon/roma/gen/api/client/message"
	"github.com/go-pantheon/roma/gen/api/client/sequence"
	"github.com/go-pantheon/roma/gen/app/codec"
	"github.com/go-pantheon/roma/gen/app/player/service"
	"github.com/pkg/errors"
)

func handleDev(ctx context.Context, s *service.PlayerServices, mod, seq int32, obj int64, in []byte) ([]byte, error) {
	cs, err := codec.UnmarshalCSDev(seq, in)
	if err != nil {
		return nil, err
	}

	var (
		sc  proto.Message
	)
	switch cliseq.DevSeq(seq) {

	// Dev command list
	case cliseq.DevSeq_DevList:
		sc, err = s.Dev.DevList(ctx, cs.(*climsg.CSDevList))

	// Execute Dev command
	case cliseq.DevSeq_DevExecute:
		sc, err = s.Dev.DevExecute(ctx, cs.(*climsg.CSDevExecute))

	default:
		return nil, errors.Errorf("seq not found. mod=%s seq=%d", "Dev", seq)
	}

	out, err0 := NewPlayerResponse(mod, seq, obj, sc)
	if err0 != nil {
		return nil, errors.Wrapf(err0, "proto marshal failed. mod=%s seq=%d", "Dev", seq)
	}
	return out, err
}
