package core

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	userobj "github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
)

var _ life.Persistent = (*UserPersister)(nil)

type UserPersister struct {
	mu sync.Mutex

	uid      int64
	user     *userobj.User
	snapshot *atomic.Value // TODO encode for other apps

	do *domain.UserDomain
}

func newUserPersister(ctx context.Context, do *domain.UserDomain, uid int64, allowCreate bool) (ret life.Persistent, newborn bool, err error) {
	p, newborn, err := do.TakeUserProto(ctx, uid, allowCreate)
	if err != nil {
		return
	}

	defer dbv1.UserProtoPool.Put(p)

	user := userobj.NewUser(uid, p.ServerVersion)
	if err = user.DecodeServer(p); err != nil {
		return
	}

	user.SetNewborn(newborn)

	ret = &UserPersister{
		uid:      uid,
		user:     user,
		do:       do,
		snapshot: &atomic.Value{},
	}

	return
}

func (s *UserPersister) Refresh(ctx context.Context) (err error) {
	s.encodeSnapshot()
	return
}

func (s *UserPersister) PrepareToPersist(ctx context.Context, modules []life.ModuleKey) (ret life.VersionProto, err error) {
	err = s.Lock(func() error {
		s.user.Version += 1 // update version first

		p := dbv1.UserProtoPool.Get()
		s.user.EncodeServer(p, modules)

		ret = p

		return nil
	})

	return
}

func (s *UserPersister) Persist(ctx context.Context, uid int64, proto life.VersionProto) (err error) {
	p, ok := proto.(*dbv1.UserProto)
	if !ok {
		err = errors.Wrapf(xerrors.ErrDBRecordType, "uid=%d proto=%T", uid, proto)
		return
	}

	return s.do.Persist(ctx, uid, p)
}

func (s *UserPersister) IncVersion(ctx context.Context, uid int64, newVersion int64) (err error) {
	return s.do.IncVersion(ctx, uid, newVersion)
}

func (s *UserPersister) OnStop(ctx context.Context, id int64, p life.VersionProto) (err error) {
	s.do.OnLogout(ctx, s.uid, p.(*dbv1.UserProto))
	return nil
}

func (s *UserPersister) Snapshot() life.VersionProto {
	return s.snapshot.Load().(life.VersionProto)
}

func (s *UserPersister) encodeSnapshot() {
	// TODO encode snapshot
}

func (s *UserPersister) Lock(f func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return f()
}

func (s *UserPersister) ModuleKeys() []life.ModuleKey {
	return userregister.AllModuleKeys()
}

func (s *UserPersister) ID() int64 {
	return s.uid
}

func (s *UserPersister) Version() int64 {
	return s.user.Version
}

func (s *UserPersister) UnsafeObject() any {
	return s.user
}
