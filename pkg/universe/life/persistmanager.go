package life

import (
	"context"

	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/pkg/errors"
)

type PersistManager struct {
	log       *log.Helper
	persister Persistent

	changed         atomic.Bool
	immediatelyChan chan struct{}
	saveChan       chan VersionProto
}

func newPersistManager(log *log.Helper, persister Persistent) (s *PersistManager) {
	s = &PersistManager{
		log:       log,
		persister: persister,
	}

	s.saveChan = make(chan VersionProto, constants.WorkerHolderSize)
	s.immediatelyChan = make(chan struct{}, constants.WorkerHolderSize)
	return
}

func (s *PersistManager) Persister() Persistent {
	return s.persister
}

func (s *PersistManager) Change(ctx context.Context, imme bool) {
	s.changed.Store(true)
	s.persister.Refresh(ctx)
	if imme {
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
	if !s.changed.CompareAndSwap(true, false) {
		return nil
	}
	return s.persister.PrepareToPersist(ctx)
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

	proto := s.persister.PrepareToPersist(ctx)
	if proto == nil {
		s.log.Infof("prepare to persist failed when incr db version. oid=%d", s.ID())
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, constants.AsyncMongoTimeout)
	defer cancel()

	ver := proto.GetVersion()
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
