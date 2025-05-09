package core

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

var _ life.Persistent = (*UserPersister)(nil)

type UserPersister struct {
	mu sync.Mutex

	uid  int64
	user *userobj.User
	show *atomic.Value

	do *domain.UserDomain
}

func newUserPersister(ctx context.Context, do *domain.UserDomain, uid int64, allowCreate bool) (ret life.Persistent, newborn bool, err error) {
	p := do.LoadCachedProto(ctx, uid, time.Now())

	if p == nil {
		p, putFunc := do.GetProto()
		defer putFunc()

		p.Id = uid
		if err = do.Load(ctx, uid, p); err != nil {
			if errors.Is(err, xerrors.ErrDBRecordNotFound) {
				if allowCreate {
					err = do.Create(ctx, uid, time.Now(), p)
					newborn = true
				}
			}
		}
		if err != nil {
			return
		}
	}

	user := userobj.NewUser(p.Id, p.Name)
	if err = user.DecodeServer(ctx, p); err != nil {
		return
	}
	user.SetNewborn(newborn)

	ret = &UserPersister{
		uid:  uid,
		user: user,
		do:   do,
		show: &atomic.Value{},
	}
	// TODO encode ShowProto
	return
}

func (s *UserPersister) Refresh(ctx context.Context) (err error) {
	s.refreshProto()
	return
}

func (s *UserPersister) PrepareToPersist(ctx context.Context) (ret life.VersionProto) {
	_ = s.Lock(func() error {
		s.user.Version += 1 // update version first

		p, putFunc := s.do.GetProto()
		defer putFunc()

		s.user.EncodeServer(p)
		ret = p
		return nil
	})
	return
}

func (s *UserPersister) refreshProto() {
	// TODO encode ShowProto
}

func (s *UserPersister) Persist(ctx context.Context, uid int64, proto life.VersionProto) (err error) {
	defer s.do.PutBackProto(proto.(*dbv1.UserProto))
	return s.do.Persist(ctx, uid, proto)
}

func (s *UserPersister) IncVersion(ctx context.Context, uid int64, newVersion int64) (err error) {
	return s.do.IncVersion(ctx, uid, newVersion)
}

func (s *UserPersister) OnStop(ctx context.Context, id int64, p life.VersionProto) (err error) {
	cache, _ := s.do.GetProto() // cache will be reset on remove from UserProtoPool.cache
	proto.Merge(cache, p)

	s.do.CacheProto(ctx, s.uid, cache, time.Now())
	return nil
}

func (s *UserPersister) ID() int64 {
	return s.uid
}

func (s *UserPersister) UnsafeObject() interface{} {
	return s.user
}

func (s *UserPersister) ShowProto() proto.Message {
	return s.show.Load().(proto.Message)
}

func (s *UserPersister) Lock(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}
