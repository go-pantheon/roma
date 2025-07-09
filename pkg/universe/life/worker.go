package life

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/profile"
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
	globalWorkerUnique = &atomic.Uint64{}
)

type Workable interface {
	xsync.Stoppable
	EventManageable

	ID() int64
}

var _ Workable = (*Worker)(nil)
var _ Responsive = (*Worker)(nil)

type Worker struct {
	xsync.Stoppable
	Broadcastable
	Responsive
	*BuiltinEventFuncs
	*Tickers

	appRouteTable    routetable.ReNewalRouteTable
	nextRTRenewAt    time.Time
	rtRenewFailCount int

	log      *log.Helper
	status   intrav1.OnlineStatus
	referer  string
	clientIP string

	// unique identifier for the current worker, used to distinguish between multiple instances of the same ID, one of usages is the optimistic lock for deletion from the Manager
	unique uint64
	events chan EventFunc
	// different error that occurred during the disconnect for logout message, default is xsync.ErrStopByTrigger
	disconnectErr  atomic.Value
	persistManager *PersistManager

	// drainingEvents flag, indicating that the events are being drained, allowing new events to be added during the process
	drainingEvents atomic.Bool

	notifyStoppedFunc func(userId int64, index uint64)
	newContextFunc    newContextFunc
}

func newWorker(
	ctx context.Context,
	log *log.Helper,
	appRouteTable routetable.ReNewalRouteTable,
	replier Responsive,
	broadcaster Broadcastable,
	builtinEventFuncs *BuiltinEventFuncs,
	persistManager *PersistManager,
	notifyStoppedFunc func(uid int64, index uint64),
	newContextFunc newContextFunc,
) (w *Worker) {
	w = &Worker{
		Stoppable:         xsync.NewStopper(constants.WorkerStopTimeout),
		unique:            globalWorkerUnique.Add(1),
		log:               log,
		appRouteTable:     appRouteTable,
		Broadcastable:     broadcaster,
		Responsive:        replier,
		BuiltinEventFuncs: builtinEventFuncs,
		Tickers:           newTickers(),
		persistManager:    persistManager,
		notifyStoppedFunc: notifyStoppedFunc,
		newContextFunc:    newContextFunc,
	}

	w.disconnectErr.Store(xsync.ErrStopByTrigger)
	w.events = make(chan EventFunc, constants.WorkerEventSize)

	w.status = OnlineStatus(xcontext.Status(ctx))
	w.referer = xcontext.GateReferer(ctx)
	w.clientIP = xcontext.ClientIP(ctx)

	w.log.WithContext(ctx).Debugf("create worker. %s", w.LogInfo())

	return w
}

func (w *Worker) start(ctx context.Context) {
	w.GoAndStop(fmt.Sprintf("worker-%d", w.ID()), func() error {
		return w.run(ctx)
	}, func() error {
		return w.Stop(ctx)
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
			return w.runEventLoop(ctx)
		})
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			return w.runTunnelResponseLoop(ctx)
		})
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			return w.runPersistLoop(ctx)
		})
	})
	eg.Go(func() error {
		return xsync.Run(func() error {
			return w.runRenewalLoop(ctx)
		})
	})

	err := eg.Wait()
	if err != nil {
		w.sendLogoutMsg(ctx, err)
	}

	return err
}

func (w *Worker) runEventLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.secondTicker.C:
			if IsGateStatus(w.Status()) {
				if err := w.EmitEvent(EventTypeSecondTick); err != nil {
					return err
				}
			}
		case <-w.minuteTicker.C:
			if IsGateStatus(w.Status()) {
				if err := w.EmitEvent(EventTypeMinuteTick); err != nil {
					return err
				}
			}
		case e, ok := <-w.ConsumeEvent():
			if !ok {
				return xerrors.ErrLifeStopped
			}

			if err := w.ExecuteEvent(w.newContextFunc(ctx, w), e); err != nil {
				return err
			}
		case <-w.persistManager.Immediately():
			if err := w.persistManager.PrepareToPersist(ctx); err != nil {
				return err
			}
		case <-w.persistTicker.C:
			if err := w.persistManager.PrepareToPersist(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Worker) runTunnelResponseLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-w.ConsumeTunnelResponse():
			if !ok {
				return xerrors.ErrLifeStopped
			}

			if err := w.ExecuteSend(msg); err != nil {
				w.log.WithContext(ctx).Errorf("worker execute reply failed. i=%d seq=<%d-%d> %s %+v", msg.GetIndex(), msg.GetMod(), msg.GetSeq(), w.LogInfo(), err)
			}
		}
	}
}

func (w *Worker) runPersistLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case proto, ok := <-w.persistManager.SaveChan():
			if !ok {
				return xerrors.ErrLifeStopped
			}

			if err := w.persistManager.Persist(ctx, proto); err != nil {
				if errors.Is(err, xerrors.ErrDBRecordNotAffected) {
					return err
				} else {
					w.log.WithContext(ctx).Errorf("worker persist failed. %s %+v", w.LogInfo(), err)
				}
			}
		}
	}
}

func (w *Worker) runRenewalLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.renewalTicker.C:
			if err := w.renewalTick(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Worker) setDisconnectErr(err error) {
	w.disconnectErr.Store(err)
}

func (w *Worker) sendLogoutMsg(_ context.Context, err error) {
	msg := &climsg.SCServerLogout{}

	switch {
	case errors.Is(err, xerrors.ErrLogoutConflictingLogin), errors.Is(err, xerrors.ErrDBRecordNotAffected):
		msg.Code = climsg.SCServerLogout_ConflictingLogin
	case errors.Is(err, xerrors.ErrLogoutKickOut):
		msg.Code = climsg.SCServerLogout_AdminKickOut
	case errors.Is(err, xerrors.ErrLogoutBanned):
		msg.Code = climsg.SCServerLogout_Banned
	default:
		msg.Code = climsg.SCServerLogout_Waiting
	}

	_ = w.PushImmediately(int32(climod.ModuleID_System), int32(cliseq.SystemSeq_ServerLogout), w.ID(), msg)
}

// reuse check if the worker can be reused and update the worker status and reply function.
// if the worker is inner status or the connection is inner context, it can be reused.
// otherwise, the worker is gate context, it can be reused if the gate referer is the same.
func (w *Worker) canReuse(ctx context.Context, replier Responsive) bool {
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
		w.UpdateSendFunc(replier.SendFunc())
	}

	return true
}

func (w *Worker) Stop(ctx context.Context) (err error) {
	return w.TurnOff(func() (err error) {
		if w.notifyStoppedFunc != nil {
			w.notifyStoppedFunc(w.ID(), w.Unique())
		}

		w.Tickers.Stop()

		wctx := w.newContextFunc(ctx, w)
		if drainErr := w.drainEvents(wctx); drainErr != nil {
			err = errors.Join(err, drainErr)
		}

		if persistErr := w.persistManager.Stop(ctx); persistErr != nil {
			err = errors.Join(err, persistErr)
		}

		return err
	})
}

func (w *Worker) EmitEvent(t WorkerEventType, args ...WithArg) error {
	if w.OnStopping() {
		return xerrors.ErrLifeStopped
	}

	f, err := w.eventFunc(t, args...)
	if err != nil {
		return err
	}

	return w.EmitEventFunc(f)
}

func (w *Worker) EmitEventFunc(f EventFunc) error {
	// if the events are being drained, allow new events to be added
	if w.drainingEvents.Load() {
		select {
		case w.events <- f:
			return nil
		default:
			return errors.Errorf("worker event channel is full on draining. %s", w.LogInfo())
		}
	}

	if w.OnStopping() {
		return xerrors.ErrLifeStopped
	}

	w.events <- f

	return nil
}

func (w *Worker) eventFunc(t WorkerEventType, args ...WithArg) (f EventFunc, e error) {
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

		if ffs := customEventFuncMap.get(t); len(ffs) > 0 {
			f = func(wctx Context) (err error) {
				for _, ff := range ffs {
					if execErr := ff(wctx, arg); execErr != nil {
						err = errors.Join(err, execErr)
					}
				}

				return nil
			}

			return f, nil
		}
	}

	if f == nil {
		return nil, errors.Errorf("worker prepared event func not found. type=%d %s", t, w.LogInfo())
	}

	return f, nil
}

func (w *Worker) ConsumeEvent() <-chan EventFunc {
	return w.events
}

func (w *Worker) ExecuteEvent(wctx Context, f EventFunc) error {
	return w.persistManager.Persister().Lock(func() error {
		if err := f(wctx); err != nil {
			return err
		}

		mods, immediately := wctx.ChangedModules()
		if len(mods) > 0 {
			w.persistManager.Change(wctx, mods, immediately)
		}

		return nil
	})
}

// drainEvents drain the events safely, including the new events generated during the event processing
func (w *Worker) drainEvents(wctx Context) (err error) {
	w.drainingEvents.Store(true)
	defer w.drainingEvents.Store(false)

	var (
		consecutiveEmptyChecks int
		maxEmptyChecks         = 3 // if the consecutive empty checks is 3, the drain is considered to be completed
		checkInterval          = time.Millisecond * 10
	)

	for {
		hasEvent := false

		select {
		case e, ok := <-w.ConsumeEvent():
			if !ok {
				break
			}

			hasEvent = true

			if executeErr := w.ExecuteEvent(wctx, e); executeErr != nil {
				err = errors.Join(err, executeErr)
			}
		default:
		}

		if hasEvent {
			consecutiveEmptyChecks = 0
		} else {
			consecutiveEmptyChecks++
			if consecutiveEmptyChecks >= maxEmptyChecks {
				break
			}

			time.Sleep(checkInterval)
		}
	}

	close(w.events)

	// handle the remaining events after the channel is closed
	for e := range w.ConsumeEvent() {
		if executeErr := w.ExecuteEvent(wctx, e); executeErr != nil {
			err = errors.Join(err, executeErr)
		}
	}

	return err
}

func (w *Worker) renewalTick(ctx context.Context) error {
	if ct := time.Now(); ct.After(w.nextRTRenewAt) {
		w.nextRTRenewAt = ct.Add(w.appRouteTable.TTL() / 2)

		if err := w.appRouteTable.RenewSelf(ctx, profile.Color(), w.ID(), profile.GRPCEndpoint()); err != nil {
			if errors.Is(err, xerrors.ErrRouteTableNotFound) || errors.Is(err, xerrors.ErrRouteTableValueNotSame) {
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

func (w *Worker) Unique() uint64 {
	return w.unique
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
