package life

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/pkg/errors"
)

type PersistManager struct {
	log       *log.Helper
	persister Persistent

	changedModules *ChangedModules

	immediatelyChan chan struct{}
	saveChan        chan VersionProto
}

type ChangedModules struct {
	sync.Mutex

	modules map[ModuleKey]struct{}
}

func (c *ChangedModules) Add(modules []ModuleKey) {
	c.Lock()
	defer c.Unlock()

	for _, module := range modules {
		c.modules[module] = struct{}{}
	}
}

func (c *ChangedModules) Take() []ModuleKey {
	c.Lock()
	defer c.Unlock()

	ret := make([]ModuleKey, 0, len(c.modules))
	for module := range c.modules {
		ret = append(ret, module)
	}

	c.modules = make(map[ModuleKey]struct{}, len(ret))

	return ret
}

func newPersistManager(log *log.Helper, persister Persistent) (s *PersistManager) {
	s = &PersistManager{
		log:       log,
		persister: persister,
		changedModules: &ChangedModules{
			modules: make(map[ModuleKey]struct{}),
		},
	}

	s.saveChan = make(chan VersionProto, constants.WorkerHolderSize)
	s.immediatelyChan = make(chan struct{}, constants.WorkerHolderSize)
	return
}

func (s *PersistManager) Persister() Persistent {
	return s.persister
}

func (s *PersistManager) Change(ctx context.Context, modules []ModuleKey, immediately bool) {
	s.changedModules.Add(modules)
	s.persister.Refresh(ctx)

	if immediately {
		s.immediatelyChan <- struct{}{}
	}
}

func (s *PersistManager) SaveChan() chan VersionProto {
	return s.saveChan
}

func (s *PersistManager) Stop(ctx context.Context) {
	close(s.immediatelyChan)
	close(s.saveChan)

	// Persist the object before stopping
	proto := s.prepareToPersist(ctx)
	if proto == nil {
		s.log.Infof("prepare to persist failed. oid=%d", s.ID())
		return
	}

	_ = s.persister.OnStop(ctx, s.ID(), proto)

	// proto will be reset by persister, so we should use it before persist, such as OnStop
	if err := s.Persist(ctx, proto); err != nil {
		s.log.WithContext(ctx).Errorf("%+v", err)
		return
	}
	s.log.WithContext(ctx).Debugf("persist stopped. oid=%d", s.ID())
}

func (s *PersistManager) Immediately() chan struct{} {
	return s.immediatelyChan
}

func (s *PersistManager) PrepareToPersist(ctx context.Context) {
	proto := s.prepareToPersist(ctx)
	if proto == nil {
		s.log.Infof("prepare to persist failed. oid=%d", s.ID())
		return
	}
	s.saveChan <- proto
}

func (s *PersistManager) prepareToPersist(ctx context.Context) VersionProto {
	if IsAdminID(s.ID()) {
		return nil
	}

	modules := s.changedModules.Take()
	if len(modules) == 0 {
		return nil
	}

	return s.persister.PrepareToPersist(ctx, modules)
}

func (s *PersistManager) Persist(c context.Context, proto VersionProto) error {
	if IsAdminID(s.ID()) {
		return nil
	}

	ctx, cancel := context.WithTimeout(c, constants.AsyncMongoTimeout)
	defer cancel()

	if err := s.persister.Persist(ctx, s.ID(), proto); err != nil {
		if errors.Is(err, xerrors.ErrDBRecordNotAffected) {
			return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "oid=%d version=%d", s.ID(), proto.GetVersion())
		}
		return errors.Errorf("persist object failed. oid=%d version=%d", s.ID(), proto.GetVersion())
	}

	if profile.IsDev() {
		s.log.WithContext(ctx).Debugf("persist object succeeded. oid=%d version=%d", s.ID(), proto.GetVersion())
	}
	return nil
}

func (s *PersistManager) IncVersion(ctx context.Context) error {
	if IsAdminID(s.ID()) {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, constants.AsyncMongoTimeout)
	defer cancel()

	ver := s.persister.Version()
	if err := s.persister.IncVersion(ctx, s.ID(), ver); err != nil {
		if errors.Is(err, xerrors.ErrDBRecordNotAffected) {
			return errors.Wrapf(err, "version error, reload. oid=%d version=%d", s.ID(), ver)
		}
		return errors.Errorf("incr db version failed. oid=%d version=%d", s.ID(), ver)
	}

	s.log.WithContext(ctx).Debugf("incr db version succeeded. oid=%d", s.ID())
	return nil
}

func (s *PersistManager) ID() int64 {
	return s.persister.ID()
}
