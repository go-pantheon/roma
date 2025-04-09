package core

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	roomobj "github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var _ life.Persistent = (*RoomPersister)(nil)

type RoomPersister struct {
	mu sync.Mutex

	rid  int64
	room *roomobj.Room
	show *atomic.Value

	do *domain.RoomDomain
}

func newRoomPersister(ctx context.Context, do *domain.RoomDomain, rid int64, allowCreate bool) (ret life.Persistent, newborn bool, err error) {
	p := do.OfflineCache(ctx, rid, time.Now())
	if p == nil {
		p := do.GetProtoFromPool()
		defer do.PutBackProtoIntoPool(p)

		p.Id = rid
		if err = do.Load(ctx, rid, p); err != nil {
			if errors.Is(err, xerrors.ErrDBRecordNotFound) {
				if allowCreate {
					err = do.Create(ctx, rid, time.Now(), p)
					newborn = true
				}
			}
		}
		if err != nil {
			return
		}
	}

	room := roomobj.NewRoom()
	if err = room.DecodeServer(p); err != nil {
		return
	}

	persister := &RoomPersister{
		rid:  p.Id,
		room: room,
		do:   do,
		show: &atomic.Value{},
	}
	return persister, newborn, nil
}

func (s *RoomPersister) Refresh(ctx context.Context) (err error) {
	s.refreshProto()
	return
}

func (s *RoomPersister) PrepareToPersist(ctx context.Context) (vp life.VersionProto) {
	_ = s.Lock(func() error {
		s.room.Version += 1          // update version first
		p := s.do.GetProtoFromPool() // p is get from sync.Pool, and will be reset by Persist() soon
		s.room.EncodeServer(p)
		vp = p
		return nil
	})
	return
}

func (s *RoomPersister) refreshProto() {
	// p := s.do.CreateProto()
	// s.room.EncodeServer(p)
	// s.proto = p
}

func (s *RoomPersister) Persist(ctx context.Context, uid int64, proto life.VersionProto) (err error) {
	defer s.do.PutBackProtoIntoPool(proto.(*dbv1.RoomProto))
	return s.do.Persist(ctx, uid, proto)
}

func (s *RoomPersister) IncVersion(ctx context.Context, uid int64, newVersion int64) (err error) {
	return s.do.IncVersion(ctx, uid, newVersion)
}

func (s *RoomPersister) OnStop(ctx context.Context, id int64, p life.VersionProto) (err error) {
	cache := s.do.GetProtoFromPool()
	proto.Merge(cache, p)

	s.do.UpdateOfflineCache(ctx, id, cache, time.Now())
	return nil
}

func (s *RoomPersister) ID() int64 {
	return s.rid
}

func (s *RoomPersister) UnsafeObject() interface{} {
	return s.room
}

func (s *RoomPersister) ShowProto() proto.Message {
	return s.show.Load().(proto.Message)
}

func (s *RoomPersister) Lock(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}
