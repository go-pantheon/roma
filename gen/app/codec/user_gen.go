// Code generated by gen-api. DO NOT EDIT.

package codec

import (
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
)

func UnmarshalCSUser(seq int32, data []byte) (p proto.Message, err error) {
	switch cliseq.UserSeq(seq) {
	// Login
	case cliseq.UserSeq_Login:
		pp := &climsg.CSLogin{}
		err = proto.Unmarshal(data, pp)
		p = pp
	// Update name
	case cliseq.UserSeq_UpdateName:
		pp := &climsg.CSUpdateName{}
		err = proto.Unmarshal(data, pp)
		p = pp
	// Set gender
	case cliseq.UserSeq_SetGender:
		pp := &climsg.CSSetGender{}
		err = proto.Unmarshal(data, pp)
		p = pp

	default:
		err = errors.Errorf("Unmarshal CSUser faild. sequence not found. seq=%d", seq)
		return
	}

	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshal CSUser faild. seq=%d", seq)
	}
	return
}

func UnmarshalSCUser(seq int32, data []byte) (p proto.Message, err error) {
	switch cliseq.UserSeq(seq) {
	// Login
	case cliseq.UserSeq_Login:
		pp := &climsg.SCLogin{}
		err = proto.Unmarshal(data, pp)
		p = pp
	// @push latest user data. The client receives the data and updates its own data to avoid data inconsistency between the client and the server when GM modifies data or the server restarts
	case cliseq.UserSeq_PushSyncUser:
		pp := &climsg.SCPushSyncUser{}
		err = proto.Unmarshal(data, pp)
		p = pp
	// Update name
	case cliseq.UserSeq_UpdateName:
		pp := &climsg.SCUpdateName{}
		err = proto.Unmarshal(data, pp)
		p = pp
	// Set gender
	case cliseq.UserSeq_SetGender:
		pp := &climsg.SCSetGender{}
		err = proto.Unmarshal(data, pp)
		p = pp

	default:
		err = errors.Errorf("Unmarshal SCUser faild. sequence not found. seq=%d", seq)
		return
	}

	if err != nil {
		return nil, errors.Wrapf(err, "Unmarshal SCUser faild. seq=%d", seq)
	}
	return
}

func IsPushSCUser(seq int32) bool {
	name := cliseq.UserSeq_name[seq]
	return strings.Index(name, "Push_") == 0
}
