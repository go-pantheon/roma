package core

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	roomobj "github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

var _ life.Persistent = (*RoomPersister)(nil)

type RoomPersister struct {
	mu sync.Mutex

	rid      int64
	room     *roomobj.Room
	snapshot *atomic.Value

	do *domain.RoomDomain
}

func newRoomPersister(ctx context.Context, do *domain.RoomDomain, oid int64, allowCreate bool) (ret life.Persistent, newborn bool, err error) {
	p, newborn, err := do.TakeRoomProto(ctx, oid, allowCreate)
	if err != nil {
		return
	}

	defer dbv1.RoomProtoPool.Put(p)

	room := roomobj.NewRoom()
	if err = room.DecodeServer(p); err != nil {
		return
	}

	persister := &RoomPersister{
		rid:      p.Id,
		room:     room,
		do:       do,
		snapshot: &atomic.Value{},
	}

	return persister, newborn, nil
}

func (s *RoomPersister) Refresh(ctx context.Context) (err error) {
	s.encodeSnapshot()
	return
}

func (s *RoomPersister) PrepareToPersist(ctx context.Context, keys []life.ModuleKey) (ret life.VersionProto, err error) {
	err = s.Lock(func() error {
		// update version first
		s.room.Version += 1

		p := dbv1.RoomProtoPool.Get()
		s.room.EncodeServer(p, keys)

		ret = p

		return nil
	})

	return ret, err
}

func (s *RoomPersister) encodeSnapshot() {
}

func (s *RoomPersister) Persist(ctx context.Context, uid int64, proto life.VersionProto) (err error) {
	defer dbv1.RoomProtoPool.Put(proto.(*dbv1.RoomProto))

	return s.do.Persist(ctx, uid, proto)
}

func (s *RoomPersister) IncVersion(ctx context.Context, uid int64, newVersion int64) (err error) {
	return s.do.IncVersion(ctx, uid, newVersion)
}

func (s *RoomPersister) OnStop(ctx context.Context, id int64, p life.VersionProto) (err error) {
	return s.do.OnLogout(ctx, id, p.(*dbv1.RoomProto))
}

func (s *RoomPersister) Snapshot() life.VersionProto {
	return s.snapshot.Load().(life.VersionProto)
}

func (s *RoomPersister) Lock(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

func (s *RoomPersister) ID() int64 {
	return s.rid
}

func (s *RoomPersister) UnsafeObject() any {
	return s.room
}

func (s *RoomPersister) Version() int64 {
	return s.room.Version
}

func (s *RoomPersister) ModuleKeys() []life.ModuleKey {
	return []life.ModuleKey{}
}
