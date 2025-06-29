package life

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Persistent interface {
	ObjectHolder

	ModuleKeys() []ModuleKey
	Lock(f func() error) error
	Refresh(ctx context.Context) (err error)
	PrepareToPersist(ctx context.Context, modules []ModuleKey) (VersionProto, error)
	Persist(ctx context.Context, id int64, proto VersionProto) (err error)
	IncVersion(ctx context.Context, id int64) (err error)
	OnStop(ctx context.Context, id int64, proto VersionProto) (err error)
}

type ObjectHolder interface {
	ID() int64
	Version() int64
	UnsafeObject() any
	Snapshot() VersionProto
}

type VersionProto interface {
	proto.Message

	GetVersion() int64
}
