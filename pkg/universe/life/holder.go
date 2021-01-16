package life

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Persistent interface {
	ID() int64
	UnsafeObject() interface{}
	ShowProto() proto.Message
	Lock(f func() error) error
	Refresh(ctx context.Context) (err error)
	PrepareToPersist(ctx context.Context) VersionProto
	Persist(ctx context.Context, id int64, proto VersionProto) (err error)
	IncVersion(ctx context.Context, id int64, newVersion int64) (err error)
	OnStop(ctx context.Context, id int64, proto VersionProto) (err error)
}

type VersionProto interface {
	Versionable
	proto.Message
}

type Versionable interface {
	GetVersion() int64
}
