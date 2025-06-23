package life

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

type Manager struct {
	xsync.Stoppable
	*PreparedTickFuncs

	log *log.Helper

	group             singleflight.Group
	statisticTicker   *time.Ticker
	workers           *WorkerMap
	stoppedWorkerChan chan string
	appRouteTable     routetable.ReNewalRouteTable

	newContext   newContextFunc
	newPersister newPersisterFunc

	approximateOnlineCount int
}

type newContextFunc func(ctx context.Context, w *Worker) Context
type newPersisterFunc func(ctx context.Context, id int64, sid int64, allowBorn bool) (persister Persistent, born bool, err error)

func NewManager(logger log.Logger, rt routetable.ReNewalRouteTable, newContext newContextFunc, newPersister newPersisterFunc) (m *Manager, stopFunc func() error) {
	m = &Manager{
		PreparedTickFuncs: newPreparedTickFuncs(),
		log:               log.NewHelper(log.With(logger, "module", "universe/life/manager")),
		appRouteTable:     rt,
		newContext:        newContext,
		newPersister:      newPersister,
	}

	m.Stoppable = xsync.NewStopper(10 * time.Second)
	m.statisticTicker = time.NewTicker(constants.ManagerStatisticTickDuration)
	m.stoppedWorkerChan = make(chan string, constants.WorkerSize)
	m.workers = NewWorkerMap()

	m.GoAndQuickStop("life-manager", func() error {
		return m.run()
	}, func() error {
		return m.Stop(context.Background())
	}, errs.IsUnloggableErr)

	return m, func() error {
		return m.Stop(context.Background())
	}
}

func (m *Manager) Worker(ctx context.Context, oid int64, sid int64, replier Replier, broadcaster Broadcaster) (worker *Worker, err error) {
	v, err, _ := m.group.Do(workerSingleFlightKey(oid), func() (any, error) {
		status := OnlineStatus(xcontext.Status(ctx))

		if old := m.get(ctx, oid); old != nil {
			if old.canReuse(ctx, replier) {
				return old, nil
			}

			newReferer := xcontext.GateReferer(ctx)
			newClientIP := xcontext.ClientIP(ctx)

			m.log.WithContext(ctx).Debugf("worker exists, stop it and replace by new worker. oid=%d old-referer=%s old-clientIP=%s new-referer=%s new-clientIP=%s", oid, old.Referer(), old.ClientIP(), newReferer, newClientIP)
			old.setDisconnectErr(xerrors.ErrLogoutConflictingLogin)

			m.workers.Remove(oid)

			if err := old.Stop(ctx); err != nil {
				m.log.WithContext(ctx).Errorf("worker stop failed. oid=%d", oid)
			}
		}

		allowBorn := !IsInnerStatus(status)

		return m.load(ctx, oid, sid, replier, broadcaster, allowBorn)
	})
	if err != nil {
		return
	}

	worker, ok := v.(*Worker)
	if !ok {
		err = errors.Errorf("not life.Worker type")
		return
	}

	return
}

func workerSingleFlightKey(id int64) string {
	return fmt.Sprintf("worker-%d", id)
}

func (m *Manager) get(_ context.Context, id int64) *Worker {
	return m.workers.Get(id)
}

func (m *Manager) load(ctx context.Context, oid int64, sid int64, replier Replier, broadcaster Broadcaster, allowBorn bool) (*Worker, error) {
	persister, born, err := m.newPersister(ctx, oid, sid, allowBorn)
	if err != nil {
		return nil, errors.WithMessagef(err, "allowBorn=%v status=%d", allowBorn, OnlineStatus(xcontext.Status(ctx)))
	}

	w := newWorker(ctx, m.log, m.appRouteTable, newPersistManager(m.log, persister),
		replier, broadcaster, newTickers(m.PreparedTickFuncs),
		m.notifyWorkerStopped, m.newContext)

	// Increment version number immediately after loading to prevent other services from loading the same data simultaneously
	// Uses version number for optimistic locking. When this worker saves data next time, the database will determine if the data is outdated through version number
	if err = w.persistManager.IncVersion(ctx); err != nil {
		return nil, err
	}

	wctx := m.newContext(ctx, w)

	if born {
		if err = m.onCreated(wctx); err != nil {
			return nil, err
		}
	}

	if err = m.onLoad(wctx); err != nil {
		return nil, err
	}

	modules, immediately := wctx.ChangedModules()
	if len(modules) > 0 {
		w.persistManager.Change(ctx, modules, immediately)
	}

	w.start(ctx)

	old := m.workers.Set(persister.ID(), w)
	if old != nil {
		m.log.WithContext(ctx).Errorf("worker already exists on load. id=%d", old.ID())

		if err := old.Stop(ctx); err != nil {
			m.log.WithContext(ctx).Errorf("worker stop failed. id=%d", old.ID())
		}
	}

	return w, nil
}

func (m *Manager) onCreated(wctx Context) (err error) {
	for _, f := range m.onCreatedEventFuncs {
		if err = f(wctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) onLoad(wctx Context) (err error) {
	for _, f := range m.onLoadEventFuncs {
		if err = f(wctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) run() error {
	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		select {
		case <-m.StopTriggered():
			return xsync.ErrStopByTrigger
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-m.statisticTicker.C:
					m.statisticTick(ctx)
				case key := <-m.stoppedWorkerChan:
					m.removeWorker(ctx, key)
				}
			}
		})
	})

	return eg.Wait()
}

func (m *Manager) Stop(ctx context.Context) (err error) {
	return m.TurnOff(func() error {
		defer close(m.stoppedWorkerChan)

		m.statisticTicker.Stop()

		for w := range m.workers.Iter() {
			if wStopErr := w.Val.Stop(ctx); wStopErr != nil {
				err = errors.JoinUnsimilar(err, wStopErr)
			}
		}

		return err
	})
}

func (m *Manager) statisticTick(ctx context.Context) {
	if c := m.workers.Count(); c != m.approximateOnlineCount {
		m.approximateOnlineCount = c
	}
}

func (m *Manager) notifyWorkerStopped(id int64, index uint64) {
	m.stoppedWorkerChan <- buildStoppedWorkerKey(id, index)
}

func buildStoppedWorkerKey(id int64, index uint64) string {
	return fmt.Sprintf("%d#%d", id, index)
}

func parseStoppedWorkerKey(key string) (id int64, index uint64, err error) {
	parts := strings.Split(key, "#")
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid key")
	}

	id, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	index, err = strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return id, index, nil
}

func (m *Manager) removeWorker(ctx context.Context, key string) {
	id, index, err := parseStoppedWorkerKey(key)
	if err != nil {
		m.log.WithContext(ctx).Errorf("life.Manager.remove failed, invalid key. key=%s", key)
		return
	}

	w := m.get(ctx, id)
	if w == nil {
		return
	}

	if w.Index() != index {
		return
	}

	if !w.OnStopping() {
		return
	}

	m.workers.Remove(id)
}
