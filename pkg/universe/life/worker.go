package life

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/router/routetable"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xsync"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/gate/intra/v1"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"golang.org/x/sync/errgroup"
)

var (
	globalWorkerIndex = &atomic.Uint64{}
)

type Workable interface {
	xsync.Stoppable
	EventManageable

	ID() int64
}

var _ Workable = (*Worker)(nil)
var _ Replier = (*Worker)(nil)

type Worker struct {
	xsync.Stoppable
	Broadcaster
	Replier
	*Tickers

	log      *log.Helper
	status   intrav1.OnlineStatus
	referer  string
	clientIP string

	// unique identifier for the current worker, used to distinguish between multiple instances of the same ID, one of usages is the optimistic lock for deletion from the Manager
	index  uint64
	events chan EventFunc
	// different error that occurred during the disconnect for logout message, default is xsync.GroupStopping
	disconnectErr  atomic.Value
	persistManager *PersistManager

	appRouteTable    routetable.ReNewalRouteTable
	nextRTRenewAt    time.Time
	rtRenewFailCount int

	notifyStoppedFunc func(userId int64, index uint64)
	newContextFunc    newContextFunc
}

func newWorker(
	ctx context.Context,
	log *log.Helper,
	appRouteTable routetable.ReNewalRouteTable,
	persistManager *PersistManager,
	replier Replier,
	broadcaster Broadcaster,
	tickers *Tickers,
	notifyStoppedFunc func(uid int64, index uint64),
	newContextFunc newContextFunc,
) (w *Worker) {
	w = &Worker{
		log:               log,
		appRouteTable:     appRouteTable,
		Broadcaster:       broadcaster,
		Replier:           replier,
		Tickers:           tickers,
		index:             globalWorkerIndex.Add(1),
		persistManager:    persistManager,
		notifyStoppedFunc: notifyStoppedFunc,
		newContextFunc:    newContextFunc,
	}

	w.Stoppable = xsync.NewStopper(10 * time.Second)
	w.disconnectErr.Store(xsync.ErrStopByTrigger)
	w.events = make(chan EventFunc, constants.WorkerEventSize)

	w.status = OnlineStatus(xcontext.Status(ctx))
	w.referer = xcontext.GateReferer(ctx)
	w.clientIP = xcontext.ClientIP(ctx)

	w.log.WithContext(ctx).Debugf("create worker. %s", w.LogInfo())

	return
}

func (w *Worker) start(ctx context.Context) {
	xsync.Go(fmt.Sprintf("worker-%d", w.ID()), func() error {
		return w.run(ctx)
	}, xerrors.IsUnlogErr)
}

func (w *Worker) run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-w.StopTriggered():
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
				case <-w.secondTicker.C:
					if IsGateStatus(w.Status()) {
						if err := w.ProductPreparedEvent(EventTypeSecondTick); err != nil {
							return err
						}
					}
				case <-w.minuteTicker.C:
					if IsGateStatus(w.Status()) {
						if err := w.ProductPreparedEvent(EventTypeMinuteTick); err != nil {
							return err
						}
					}
				case e := <-w.ConsumeEvent():
					if err := w.ExecuteEvent(w.newContextFunc(ctx, w), e); err != nil {
						w.log.WithContext(ctx).Errorf("worker execute event failed. %s %+v", w.LogInfo(), err)
					}
				case <-w.persistManager.Immediately():
					if err := w.persistManager.PrepareToPersist(ctx); err != nil {
						w.log.WithContext(ctx).Errorf("worker immediately prepare to persist failed. %s %+v", w.LogInfo(), err)
					}
				case <-w.persistTicker.C:
					if err := w.persistManager.PrepareToPersist(ctx); err != nil {
						w.log.WithContext(ctx).Errorf("worker ticker prepare to persist failed. %s %+v", w.LogInfo(), err)
					}
				}
			}
		})
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case msg := <-w.ConsumeReplyMessage():
					if err := w.ExecuteReply(msg); err != nil {
						w.log.WithContext(ctx).Errorf("worker execute reply failed. id=%d %+v", w.ID(), err)
					}
				}
			}
		})
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case proto := <-w.persistManager.SaveChan():
					if err := w.persistManager.Persist(ctx, proto); err != nil {
						if errors.Is(err, xerrors.ErrDBRecordNotAffected) {
							return err
						} else {
							w.log.WithContext(ctx).Errorf("worker persist failed. id=%d %+v", w.ID(), err)
						}
					}
				}
			}
		})
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			for {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-w.workerTicker.C:
					if err := w.workerTick(ctx); err != nil {
						return err
					}
				}
			}
		})
	})
	err := eg.Wait()
	if err != nil {
		w.sendLogoutMsg(ctx, err)
	}
	return err
}

func (w *Worker) setDisconnectErr(err error) {
	w.disconnectErr.Store(err)
}

func (w *Worker) sendLogoutMsg(_ context.Context, err error) {
	if xerrors.IsUnlogErr(err) {
		return
	}

	msg := &climsg.SCServerLogout{}

	switch {
	case errors.Is(err, xerrors.ErrLogoutConflictingLogin):
		msg.Code = climsg.SCServerLogout_ConflictingLogin
	case errors.Is(err, xerrors.ErrLogoutKickOut):
		msg.Code = climsg.SCServerLogout_AdminKickOut
	case errors.Is(err, xerrors.ErrLogoutBanned):
		msg.Code = climsg.SCServerLogout_Banned
	default:
		msg.Code = climsg.SCServerLogout_Waiting
	}

	_ = w.ReplyImmediately(climod.ModuleID_System, int32(cliseq.SystemSeq_ServerLogout), w.ID(), msg)
}

// reuse check if the worker can be reused and update the worker status and reply function.
// if the worker is inner status or the connection is inner context, it can be reused.
// otherwise, the worker is gate context, it can be reused if the gate referer is the same.
func (w *Worker) canReuse(ctx context.Context, replier Replier) bool {
	if !IsInnerStatus(w.Status()) && !IsInnerContext(ctx) {
		return false
	}
	if xcontext.GateReferer(ctx) != w.Referer() {
		return false
	}

	if xcontext.Status(ctx) == int64(intrav1.OnlineStatus_ONLINE_STATUS_GATE) {
		w.status = intrav1.OnlineStatus_ONLINE_STATUS_GATE
		w.referer = xcontext.GateReferer(ctx)
		w.clientIP = xcontext.ClientIP(ctx)
		w.Replier.UpdateReplyFunc(replier.GetReplyFunc())
	}
	return true
}

func (w *Worker) Stop(ctx context.Context) (err error) {
	return w.TurnOff(func() (err error) {
		w.log.Debugf("worker is stopping. id=%d", w.ID())
		if w.notifyStoppedFunc != nil {
			w.notifyStoppedFunc(w.ID(), w.Index())
		}

		w.Tickers.stop()

		close(w.events)
		wctx := w.newContextFunc(ctx, w)
		for e := range w.ConsumeEvent() {
			if executeErr := w.ExecuteEvent(wctx, e); executeErr != nil {
				err = errors.Join(err, executeErr)
			}
		}

		w.persistManager.Stop(ctx)
		w.log.Debugf("worker stopped. %s", w.LogInfo())

		return err
	})
}

func (w *Worker) ProductFuncEvent(f EventFunc) error {
	if w.OnStopping() {
		return xerrors.ErrLifeStopped
	}

	w.events <- f
	return nil
}

func (w *Worker) ProductPreparedEvent(t WorkerEventType, args ...WithArg) error {
	if w.OnStopping() {
		return xerrors.ErrLifeStopped
	}

	f, err := w.preparedEventFunc(t, args...)
	if err != nil {
		return err
	}

	w.events <- f
	return nil
}

func (w *Worker) preparedEventFunc(t WorkerEventType, args ...WithArg) (f EventFunc, e error) {
	switch t {
	case EventTypeSecondTick:
		f = w.secondTick
	case EventTypeMinuteTick:
		f = w.minuteTick
	default:
		arg := GetEventArg()
		defer PutEventArg(arg)

		for _, a := range args {
			a(arg)
		}

		if ffs, ok := preparedEventFuncMap.get(t); ok {
			f = func(wctx Context) (err error) {
				for _, ff := range ffs {
					if err := ff(wctx, arg); err != nil {
						w.log.WithContext(wctx).Errorf("worker execute event failed. id=%d type=%d %+v", wctx.OID(), t, err)
					}
				}
				return
			}
			return
		}
	}

	if f == nil {
		e = errors.Errorf("worker prepared event func not found. id=%d type=%d", w.ID(), t)
	}
	return
}

func (w *Worker) ConsumeEvent() <-chan EventFunc {
	return w.events
}

func (w *Worker) ExecuteEvent(wctx Context, f EventFunc) error {
	e := w.persistManager.Persister().Lock(func() error {
		err := f(wctx)
		mods, immediately := wctx.ChangedModules()
		if len(mods) > 0 {
			w.persistManager.Change(wctx, mods, immediately)
		}
		return err
	})

	return e
}

func (w *Worker) workerTick(ctx context.Context) error {
	if ct := time.Now(); ct.After(w.nextRTRenewAt) {
		w.nextRTRenewAt = ct.Add(w.appRouteTable.TTL() / 2)

		if err := w.appRouteTable.RenewSelf(ctx, "gate", w.ID(), w.Referer()); err != nil {
			if errors.Is(err, routetable.ErrRouteTableNotFound) || errors.Is(err, routetable.ErrRouteTableValueNotSame) {
				return err
			}

			w.rtRenewFailCount++
			w.nextRTRenewAt = ct.Add(w.appRouteTable.TTL() / 10)

			if w.rtRenewFailCount > 3 {
				return err
			}
		}
	} else {
		w.rtRenewFailCount = 0
	}

	return nil
}

func (w *Worker) ID() int64 {
	return w.persistManager.ID()
}

func (w *Worker) IsAdminID() bool {
	return IsAdminID(w.ID())
}

func IsAdminID(id int64) bool {
	return id == 0
}

func (w *Worker) Index() uint64 {
	return w.index
}

func (w *Worker) Status() intrav1.OnlineStatus {
	return w.status
}

func (w *Worker) Referer() string {
	return w.referer
}

func (w *Worker) ClientIP() string {
	return w.clientIP
}

func (w *Worker) LogInfo() string {
	buf := strings.Builder{}

	buf.WriteString("id=")
	buf.WriteString(strconv.FormatInt(w.ID(), 10))
	buf.WriteString(" status=")
	buf.WriteString(w.Status().String())
	buf.WriteString(" referer=")
	buf.WriteString(w.Referer())
	buf.WriteString(" client_ip=")
	buf.WriteString(w.ClientIP())

	return buf.String()
}
