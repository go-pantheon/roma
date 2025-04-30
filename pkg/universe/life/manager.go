package life

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/pkg/errors"
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

	newContext   newContextFunc
	newPersister newPersisterFunc

	onlineCountDev int
}

type newContextFunc func(ctx context.Context, w *Worker) Context
type newPersisterFunc func(ctx context.Context, id int64, allowBorn bool) (persister Persistent, born bool, err error)

func NewManager(logger log.Logger, newContext newContextFunc, newPersister newPersisterFunc) (m *Manager, stopFunc func()) {
	m = &Manager{
		PreparedTickFuncs: newPreparedTickFuncs(),
		log:               log.NewHelper(log.With(logger, "module", "universe/life/manager")),
		newContext:        newContext,
		newPersister:      newPersister,
	}

	m.Stoppable = xsync.NewStopper(10*time.Second, xsync.WithLog(m.log))
	m.statisticTicker = time.NewTicker(constants.ManagerStatisticTickDuration)
	m.stoppedWorkerChan = make(chan string, constants.WorkerSize)
	m.workers = NewWorkerMap()

	xsync.GoSafe("life-manager", func() error {
		return m.run()
	}, errs.DontLog)

	return m, func() {
		m.TriggerStop()
	}
}

func (m *Manager) Worker(ctx context.Context, oid int64, replier Replier, broadcaster Broadcaster) (worker *Worker, err error) {
	v, err, _ := m.group.Do(workerSingleFlightKey(oid), func() (interface{}, error) {
		status := OnlineStatus(xcontext.Status(ctx))
		if old := m.get(ctx, oid); old != nil {
			if old.reuse(ctx, replier) {
				return old, nil
			}

			newReferer := xcontext.GateReferer(ctx)
			newClientIP := xcontext.ClientIP(ctx)

			old.log.WithContext(ctx).Debugf("worker exists, stop it and replace by new worker. oid=%d old-referer=%s old-clientIP=%s new-referer=%s new-clientIP=%s", oid, old.Referer(), old.ClientIP(), newReferer, newClientIP)
			old.setDisconnectErr(xerrors.ErrLogoutConflictingLogin)
			old.TriggerStop()
			old.WaitStopped()
		}

		allowBorn := !IsInnerStatus(status)
		return m.load(ctx, oid, replier, broadcaster, allowBorn)
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

func (m *Manager) load(ctx context.Context, oid int64, replier Replier, broadcaster Broadcaster, allowBorn bool) (*Worker, error) {
	var (
		changed bool
		err     error
	)

	persister, born, err := m.newPersister(ctx, oid, allowBorn)
	if err != nil {
		return nil, errors.WithMessagef(err, "allowBorn=%v status=%d", allowBorn, OnlineStatus(xcontext.Status(ctx)))
	}

	w := newWorker(ctx, m.log, newPersistManager(m.log, persister),
		replier, broadcaster, newTickers(m.PreparedTickFuncs),
		m.notifyWorkerStopped, m.newContext)

	// Increment version number immediately after loading to prevent other services from loading the same data simultaneously
	// Uses version number for optimistic locking. When this worker saves data next time, the database will determine if the data is outdated through version number
	if err = w.persistManager.IncVersion(ctx); err != nil {
		return nil, err
	}

	wctx := m.newContext(ctx, w)

	if born {
		if changed, err = m.onCreated(wctx); err != nil {
			return nil, err
		}
	}
	if err = m.onLoad(wctx); err != nil {
		return nil, err
	}

	if changed {
		w.persistManager.Change(ctx, false)
	}

	w.start(ctx)
	old := m.workers.Set(persister.ID(), w)
	if old != nil {
		m.log.WithContext(ctx).Errorf("worker exists, stop old worker. id=%d", old.ID())
		old.TriggerStop()
	}
	return w, nil
}

func (m *Manager) onCreated(wctx Context) (changed bool, err error) {
	for _, f := range m.onCreatedEventFuncs {
		f(wctx)
	}
	changed = true
	return
}

func (m *Manager) onLoad(wctx Context) (err error) {
	for _, f := range m.onLoadEventFuncs {
		f(wctx)
	}
	return
}

func (m *Manager) run() error {
	defer func() {
		m.stop()
	}()

	eg, ctx := errgroup.WithContext(context.Background())
	eg.Go(func() error {
		select {
		case <-m.StopTriggered():
			return xsync.ErrGroupStopping
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		return xsync.RunSafe(func() error {
			for {
				select {
				case <-m.statisticTicker.C:
					m.statisticTick(ctx)
				case key := <-m.stoppedWorkerChan:
					m.removeWorker(ctx, key)
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	})

	return eg.Wait()
}

func (m *Manager) stop() {
	m.DoStop(func() {
		defer close(m.stoppedWorkerChan)
		m.statisticTicker.Stop()
		for w := range m.workers.Iter() {
			w.Val.TriggerStop()
		}
		m.log.Infof("life.Manager stopped")
	})
}

func (m *Manager) statisticTick(ctx context.Context) {
	if c := m.workers.Count(); c != m.onlineCountDev {
		m.onlineCountDev = c
		m.log.WithContext(ctx).Debugf("[stat] current online count: %d", c)
	}
}

// notifyWorkerStopped to the Manager that the Worker is about to be deleted
// Parameters:
//   - id: The ID of the Worker
//   - nonce: The Nonce of the Worker
func (m *Manager) notifyWorkerStopped(id int64, nonce string) {
	m.stoppedWorkerChan <- buildStoppedWorkerKey(id, nonce)
}

func buildStoppedWorkerKey(id int64, nonce string) string {
	return fmt.Sprintf("%d#%s", id, nonce)
}

func parseStoppedWorkerKey(key string) (id int64, nonce string, err error) {
	parts := strings.Split(key, "#")
	if len(parts) != 2 {
		return 0, "", errors.New("invalid key")
	}
	id, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", err
	}
	return id, parts[1], nil
}

func (m *Manager) removeWorker(ctx context.Context, key string) {
	id, nonce, err := parseStoppedWorkerKey(key)
	if err != nil {
		m.log.WithContext(ctx).Errorf("life.Manager.remove failed, invalid key. key=%s", key)
		return
	}

	w := m.get(ctx, id)
	if w == nil {
		m.log.WithContext(ctx).Errorf("life.Manager.remove failed, worker not found. oid=%d", id)
		return
	}
	if w.Nonce() != nonce {
		m.log.WithContext(ctx).Debugf("life.Manager.remove failed, nonce not match. oid=%d nonce=%s", id, nonce)
		return
	}
	if !w.IsStopping() {
		m.log.WithContext(ctx).Errorf("life.Manager.remove failed, worker not stopped. oid=%d", id)
		return
	}

	m.workers.Remove(id)
	m.log.WithContext(ctx).Debugf("life.Manager.remove success, oid=%d", id)
}
