package life

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Persistent interface {
	ID() int64
	Version() int64
	UnsafeObject() interface{}
	ShowProto() proto.Message
	Lock(f func() error) error
	Refresh(ctx context.Context) (err error)
	PrepareToPersist(ctx context.Context, modules []ModuleKey) VersionProto
	Persist(ctx context.Context, id int64, proto VersionProto) (err error)
	IncVersion(ctx context.Context, id int64, newVersion int64) (err error)
	OnStop(ctx context.Context, id int64, proto VersionProto) (err error)

	// PrepareToPersistV2(ctx context.Context) PersistData
	// PersistV2(ctx context.Context, data PersistData) (err error)
}


type PersistData struct {
	ID      int64
	Version int64
	Modules map[ModuleKey]string
}

type VersionProto interface {
	Versionable
	proto.Message
}

type Versionable interface {
	GetVersion() int64
}
