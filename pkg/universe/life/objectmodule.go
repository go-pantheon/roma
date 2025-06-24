package life

import "google.golang.org/protobuf/proto"

type ModuleKey string

type NewModuleFunc func() Module

type Module interface {
	ModuleCodec

	IsLifeModule()
}

type ModuleCodec interface {
	EncodeServer() proto.Message
	DecodeServer(proto.Message) error
}
