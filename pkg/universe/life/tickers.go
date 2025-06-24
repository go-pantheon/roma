package life

import (
	"time"

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
	t.workerTicker = time.NewTicker(time.Second * 10)
	t.secondTicker = time.NewTicker(time.Second)
	t.minuteTicker = time.NewTicker(time.Minute)
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
		f(wctx)
	}
	return
}

func (tf *PreparedTickFuncs) minuteTick(wctx Context) (err error) {
	for _, f := range tf.minuteTickFuncs {
		f(wctx)
	}
	return
}
