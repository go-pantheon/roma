package life

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/go-pantheon/roma/pkg/universe/data"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

type Manager struct {
	xsync.Stoppable
	*BuiltinEventFuncs

	log *log.Helper

	group             singleflight.Group
	statisticTicker   *time.Ticker
	workers           *WorkerMap
	stoppedWorkerChan chan string
	appRouteTable     routetable.ReNewalRouteTable
	pusher            *data.PushRepo

	newContext   newContextFunc
	newPersister newPersisterFunc

	approximateOnlineCount int
}

type newContextFunc func(ctx context.Context, w *Worker) Context
type newPersisterFunc func(ctx context.Context, id int64, allowBorn bool) (persister Persistent, born bool, err error)

func NewManager(logger log.Logger, rt routetable.ReNewalRouteTable, pusher *data.PushRepo, newContext newContextFunc, newPersister newPersisterFunc) (m *Manager, stopFunc func() error) {
	m = &Manager{
		Stoppable:         xsync.NewStopper(constants.ManagerStopTimeout),
		log:               log.NewHelper(log.With(logger, "module", "universe/life/manager")),
		BuiltinEventFuncs: newBuiltinEventFuncs(),
		appRouteTable:     rt,
		pusher:            pusher,
		newContext:        newContext,
		newPersister:      newPersister,
	}

	m.statisticTicker = time.NewTicker(constants.StatisticTickDuration)
	m.stoppedWorkerChan = make(chan string, constants.WorkerSize)
	m.workers = NewWorkerMap()

	m.GoAndQuickStop("life-manager", func() error {
		return m.run()
	}, func() error {
		return m.Stop(context.Background())
	}, xerrors.IsUnlogErr)

	return m, func() error {
		return m.Stop(context.Background())
	}
}

func (m *Manager) Worker(ctx context.Context, oid int64, replier Responsive) (worker *Worker, err error) {
	v, err, _ := m.group.Do(workerSingleFlightKey(oid), func() (any, error) {
		status := OnlineStatus(xcontext.Status(ctx))

		if old := m.workers.Get(oid); old != nil {
			if old.canReuse(ctx, replier) {
				return old, nil
			}

			old.setDisconnectErr(xerrors.ErrLogoutConflictingLogin)

			if err := old.Stop(ctx); err != nil {
				m.log.WithContext(ctx).Errorf("worker stop failed. oid=%d", oid)
			}
		}

		allowBorn := !IsInnerStatus(status)

		return m.load(ctx, oid, replier, NewBroadcaster(m.pusher), allowBorn)
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

func (m *Manager) load(ctx context.Context, oid int64, replier Responsive, broadcaster Broadcastable, allowBorn bool) (*Worker, error) {
	persister, born, err := m.newPersister(ctx, oid, allowBorn)
	if err != nil {
		return nil, errors.WithMessagef(err, "allowBorn=%v status=%d", allowBorn, OnlineStatus(xcontext.Status(ctx)))
	}

	// Increment version number of optimistic lock immediately to prevent other services from loading the same data simultaneously
	if err = persister.IncVersion(ctx, oid, persister.Version()); err != nil {
		return nil, err
	}

	w := newWorker(ctx, m.log, m.appRouteTable, replier, broadcaster, m.BuiltinEventFuncs,
		newPersistManager(m.log, persister), m.notifyWorkerStopped, m.newContext)

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

	if old := m.workers.Set(oid, w); old != nil {
		if err := old.Stop(ctx); err != nil {
			return nil, errors.WithMessagef(err, "old worker stop failed")
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
					m.removeStoppedWorker(ctx, key)
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

		wg := sync.WaitGroup{}

		for w := range m.workers.Iter() {
			wg.Add(1)
			xsync.Go(fmt.Sprintf("life-manager-stop-worker-%d", w.Key), func() error {
				defer wg.Done()

				if wStopErr := w.Val.Stop(ctx); wStopErr != nil {
					err = errors.JoinUnsimilar(err, wStopErr)
				}

				return nil
			})
		}

		wg.Wait()

		return err
	})
}

func (m *Manager) statisticTick(_ context.Context) {
	if profile.IsProd() {
		return
	}

	if c := m.workers.Count(); c != m.approximateOnlineCount {
		m.approximateOnlineCount = c
	}
}

func (m *Manager) notifyWorkerStopped(id int64, index uint64) {
	m.stoppedWorkerChan <- stoppedWorkerKey(id, index)
}

func stoppedWorkerKey(id int64, index uint64) string {
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

func (m *Manager) removeStoppedWorker(ctx context.Context, key string) {
	id, index, err := parseStoppedWorkerKey(key)
	if err != nil {
		m.log.WithContext(ctx).Errorf("life.Manager.remove failed, invalid key. key=%s", key)
		return
	}

	w := m.workers.Get(id)
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

func (m *Manager) Pusher() *data.PushRepo {
	return m.pusher
}
