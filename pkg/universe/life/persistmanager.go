package life

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/pkg/universe/constants"
)

type PersistManager struct {
	xsync.Stoppable

	log       *log.Helper
	persister Persistent

	changedModules *ChangedModules

	immediatelyChan chan struct{}
	saveChan        chan VersionProto
}

func newPersistManager(log *log.Helper, persister Persistent) (s *PersistManager) {
	s = &PersistManager{
		Stoppable:       xsync.NewStopper(constants.PersistManagerStopTimeout),
		log:             log,
		persister:       persister,
		changedModules:  NewChangedModules(len(persister.ModuleKeys())),
		saveChan:        make(chan VersionProto, constants.WorkerHolderSize),
		immediatelyChan: make(chan struct{}, constants.WorkerHolderSize),
	}

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

func (s *PersistManager) Stop(ctx context.Context) (err error) {
	return s.TurnOff(func() error {
		close(s.immediatelyChan)
		close(s.saveChan)

		// Persist the object before stopping
		proto, err := s.prepareToPersist(ctx)
		if err != nil {
			return err
		}

		if proto == nil {
			return nil
		}

		if onStopErr := s.persister.OnStop(ctx, s.ID(), proto); onStopErr != nil {
			s.log.WithContext(ctx).Errorf("persister on stop failed. oid=%d %+v", s.ID(), onStopErr)
		}

		// proto will be reset by persister, other functions should use it Before Persist
		return s.persist(ctx, proto)
	})
}

func (s *PersistManager) Immediately() chan struct{} {
	return s.immediatelyChan
}

func (s *PersistManager) PrepareToPersist(ctx context.Context) error {
	if s.OnStopping() {
		return xerrors.ErrLifeStopped
	}

	proto, err := s.prepareToPersist(ctx)
	if err != nil {
		return err
	}

	if proto == nil {
		return errors.New("persist nil proto")
	}

	s.saveChan <- proto

	return nil
}

func (s *PersistManager) prepareToPersist(ctx context.Context) (VersionProto, error) {
	if IsAdminID(s.ID()) {
		return nil, nil
	}

	modules := s.changedModules.Move()
	if len(modules) == 0 {
		return nil, nil
	}

	return s.persister.PrepareToPersist(ctx, modules)
}

func (s *PersistManager) Persist(ctx context.Context, proto VersionProto) error {
	if s.OnStopping() {
		return xerrors.ErrLifeStopped
	}

	return s.persist(ctx, proto)
}

func (s *PersistManager) persist(c context.Context, proto VersionProto) error {
	if IsAdminID(s.ID()) {
		return nil
	}

	ctx, cancel := context.WithTimeout(c, constants.WorkerPersistTimeout)
	defer cancel()

	return s.persister.Persist(ctx, s.ID(), proto)
}

func (s *PersistManager) IncVersion(ctx context.Context) error {
	if IsAdminID(s.ID()) {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, constants.WorkerPersistTimeout)
	defer cancel()

	return s.persister.IncVersion(ctx, s.ID(), s.persister.Version())
}

func (s *PersistManager) ID() int64 {
	return s.persister.ID()
}

type ChangedModules struct {
	mu sync.Mutex

	modules []ModuleKey
}

func NewChangedModules(cap int) *ChangedModules {
	return &ChangedModules{
		modules: make([]ModuleKey, 0, cap),
	}
}

func (c *ChangedModules) Add(modules []ModuleKey) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.modules = append(c.modules, modules...)
}

func (c *ChangedModules) Move() []ModuleKey {
	c.mu.Lock()

	defer func() {
		c.modules = make([]ModuleKey, 0, len(c.modules))
		c.mu.Unlock()
	}()

	return c.modules
}
