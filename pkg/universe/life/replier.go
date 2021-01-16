package life

import (
	climod "github.com/vulcan-frame/vulcan-game/gen/api/client/module"
	"google.golang.org/protobuf/proto"
)

type ReplyFunc = func(msg proto.Message) error

type Replier interface {
	Reply(mod climod.ModuleID, seq int32, obj int64, sc proto.Message) error
	ReplyImmediately(mod climod.ModuleID, seq int32, obj int64, sc proto.Message) error
	ReplyBytes(mod climod.ModuleID, seq int32, obj int64, body []byte) error
	ConsumeReplyMessage() <-chan proto.Message
	ExecuteReply(msg proto.Message) error
	UpdateReplyFunc(replyFunc ReplyFunc)
	GetReplyFunc() ReplyFunc
}
