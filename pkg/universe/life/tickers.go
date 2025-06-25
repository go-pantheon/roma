package life

import (
	"time"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/pkg/universe/constants"
)

type Tickers struct {
	*PreparedTickFuncs

	workerTicker  *time.Ticker
	secondTicker  *time.Ticker
	minuteTicker  *time.Ticker
	persistTicker *time.Ticker
}

type PreparedTickFuncs struct {
	secondTickFuncs     []func(ctx Context) error
	minuteTickFuncs     []func(ctx Context) error
	onCreatedEventFuncs []func(ctx Context) error
	onLoadEventFuncs    []func(ctx Context) error
}

func newPreparedTickFuncs() *PreparedTickFuncs {
	return &PreparedTickFuncs{}
}

func newTickers(fs *PreparedTickFuncs) *Tickers {
	t := &Tickers{
		PreparedTickFuncs: fs,
	}

	t.workerTicker = time.NewTicker(constants.WorkerRenewalTickDuration)
	t.secondTicker = time.NewTicker(constants.WorkerSecondTickDuration)
	t.minuteTicker = time.NewTicker(constants.WorkerMinuteTickDuration)
	t.persistTicker = time.NewTicker(constants.WorkerPersistTickDuration)

	return t
}

func (t *Tickers) stop() {
	t.workerTicker.Stop()
	t.persistTicker.Stop()
	t.secondTicker.Stop()
	t.minuteTicker.Stop()
}

func (tf *PreparedTickFuncs) RegisterSecondTick(f func(ctx Context) error) {
	tf.secondTickFuncs = append(tf.secondTickFuncs, f)
}

func (tf *PreparedTickFuncs) RegisterMinuteTick(f func(ctx Context) error) {
	tf.minuteTickFuncs = append(tf.minuteTickFuncs, f)
}

func (tf *PreparedTickFuncs) RegisterOnLoadEvent(f func(ctx Context) error) {
	tf.onLoadEventFuncs = append(tf.onLoadEventFuncs, f)
}

func (tf *PreparedTickFuncs) RegisterOnCreatedEvent(f func(ctx Context) error) {
	tf.onCreatedEventFuncs = append(tf.onCreatedEventFuncs, f)
}

func (tf *PreparedTickFuncs) RegisterCustomEvent(e WorkerEventType, f eventFunc) {
	preparedEventFuncMap.put(e, f)
}

func (tf *PreparedTickFuncs) secondTick(wctx Context) (err error) {
	for _, f := range tf.secondTickFuncs {
		if tickErr := f(wctx); tickErr != nil {
			err = errors.Join(err, tickErr)
		}
	}

	return err
}

func (tf *PreparedTickFuncs) minuteTick(wctx Context) (err error) {
	for _, f := range tf.minuteTickFuncs {
		if tickErr := f(wctx); tickErr != nil {
			err = errors.Join(err, tickErr)
		}
	}

	return err
}
