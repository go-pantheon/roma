package life

import (
	"context"
	"fmt"
	"time"

	"sync/atomic"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xcontext"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xsync"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	climod "github.com/go-pantheon/roma/gen/api/client/module"
	cliseq "github.com/go-pantheon/roma/gen/api/client/sequence"
	intrav1 "github.com/go-pantheon/roma/gen/api/server/gate/intra/v1"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/go-pantheon/roma/pkg/universe/constants"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
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
	nonce  string
	events chan EventFunc
	// different error that occurred during the disconnect for logout message, default is xsync.GroupStopping
	disconnectErr  atomic.Value
	persistManager *PersistManager

	notifyStoppedFunc func(userId int64, nonce string)
	newContextFunc    newContextFunc
}

func newWorker(ctx context.Context, log *log.Helper,
	persistManager *PersistManager,
	replier Replier,
	broadcaster Broadcaster,
	tickers *Tickers,
	notifyStoppedFunc func(uid int64, nonce string),
	newContextFunc newContextFunc,
) (w *Worker) {
	w = &Worker{
		log:               log,
		Broadcaster:       broadcaster,
		Replier:           replier,
		Tickers:           tickers,
		nonce:             uuid.New().String(),
		persistManager:    persistManager,
		notifyStoppedFunc: notifyStoppedFunc,
		newContextFunc:    newContextFunc,
	}

	w.Stoppable = xsync.NewStopper(10*time.Second, xsync.WithLog(w.log))
	w.disconnectErr.Store(xsync.ErrGroupStopping)
	w.events = make(chan EventFunc, constants.WorkerEventSize)

	w.status = OnlineStatus(xcontext.Status(ctx))
	w.referer = xcontext.GateReferer(ctx)
	w.clientIP = xcontext.ClientIP(ctx)
	w.log.WithContext(ctx).Debugf("create worker. id=%d status=%d client_ip=%s referer=%s", w.ID(), w.status, w.clientIP, w.referer)

	return
}

func (w *Worker) start(ctx context.Context) {
	xsync.GoSafe(fmt.Sprintf("worker-%d", w.ID()), func() error {
		return w.run(ctx)
	}, errs.DontLog)
}

func (w *Worker) run(ctx context.Context) error {
	defer func() {
		w.stop(ctx)
	}()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-w.StopTriggered():
			return w.disconnectErr.Load().(error)
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	eg.Go(func() error {
		return xsync.RunSafe(func() error {
			for {
				select {
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
						w.log.WithContext(ctx).Errorf("worker execute event failed. id=%d %+v", w.ID(), err)
					}
				case <-w.persistManager.Immediately():
					w.persistManager.PrepareToPersist(ctx)
				case <-w.persistTicker.C:
					w.persistManager.PrepareToPersist(ctx)
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	})
	eg.Go(func() error {
		return xsync.RunSafe(func() error {
			for {
				select {
				case msg := <-w.ConsumeReplyMessage():
					if err := w.ExecuteReply(msg); err != nil {
						w.log.WithContext(ctx).Errorf("worker execute reply failed. id=%d %+v", w.ID(), err)
					}
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		})
	})
	eg.Go(func() error {
		return xsync.RunSafe(func() error {
			for {
				select {
				case proto := <-w.persistManager.SaveChan():
					if err := w.persistManager.Persist(ctx, proto); err != nil {
						if errors.Is(err, xerrors.ErrDBRecordNotAffected) {
							return err
						} else {
							w.log.WithContext(ctx).Errorf("worker persist failed. id=%d %+v", w.ID(), err)
						}
					}
				case <-ctx.Done():
					return ctx.Err()
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
	if errs.IsConnectionError(err) || errs.IsContextError(err) {
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
func (w *Worker) reuse(ctx context.Context, replier Replier) bool {
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

func (w *Worker) stop(ctx context.Context) {
	w.DoStop(func() {
		w.log.Debugf("worker is stopping. id=%d", w.ID())
		if w.notifyStoppedFunc != nil {
			w.notifyStoppedFunc(w.ID(), w.Nonce())
		}

		w.Tickers.stop()

		close(w.events)
		wctx := w.newContextFunc(ctx, w)
		for e := range w.ConsumeEvent() {
			if err := w.ExecuteEvent(wctx, e); err != nil {
				w.log.WithContext(ctx).Errorf("worker execute event on stop failed. id=%d %+v", w.ID(), err)
			}
		}

		w.persistManager.Stop(ctx)
		w.log.Debugf("worker stopped. id=%d", w.ID())
	})
}

func (w *Worker) ProductFuncEvent(f EventFunc) error {
	if w.IsStopping() {
		return errs.ErrLifeWorkerStopped
	}

	w.events <- f
	return nil
}

func (w *Worker) ProductPreparedEvent(t WorkerEventType, args ...int64) error {
	if w.IsStopping() {
		return errs.ErrLifeWorkerStopped
	}

	f, err := w.preparedEventFunc(t, args...)
	if err != nil {
		return err
	}

	w.events <- f
	return nil
}

func (w *Worker) preparedEventFunc(t WorkerEventType, args ...int64) (f EventFunc, e error) {
	switch t {
	case EventTypeSecondTick:
		f = w.secondTick
	case EventTypeMinuteTick:
		f = w.minuteTick
	default:
		if ffs, ok := preparedEventFuncMap.get(t); ok {
			f = func(wctx Context) (err error) {
				for _, ff := range ffs {
					if err := ff(wctx, args...); err != nil {
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
		changed, immediately := wctx.IsChanged()
		if changed {
			w.persistManager.Change(wctx, immediately)
		}
		return err
	})

	return e
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

// Nonce return the nonce of the worker
// nonce is unique identifier for the current worker, used to distinguish between multiple instances of the same ID, one of usages is the optimistic lock for deletion from the Manager
func (w *Worker) Nonce() string {
	return w.nonce
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
