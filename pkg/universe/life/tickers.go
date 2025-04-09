package life

import (
	"time"

	"github.com/go-pantheon/roma/pkg/universe/constants"
)

type Tickers struct {
	*PreparedTickFuncs

	persistTicker *time.Ticker
	secondTicker  *time.Ticker
	minuteTicker  *time.Ticker
}

type PreparedTickFuncs struct {
	secondTickFuncs     []func(ctx Context)
	minuteTickFuncs     []func(ctx Context)
	onCreatedEventFuncs []func(ctx Context)
	onLoadEventFuncs    []func(ctx Context)
}

func newPreparedTickFuncs() *PreparedTickFuncs {
	return &PreparedTickFuncs{}
}

func newTickers(fs *PreparedTickFuncs) *Tickers {
	t := &Tickers{
		PreparedTickFuncs: fs,
	}
	t.secondTicker = time.NewTicker(time.Second)
	t.minuteTicker = time.NewTicker(time.Minute)
	t.persistTicker = time.NewTicker(constants.WorkerPersistTickDuration)
	return t
}

func (t *Tickers) stop() {
	t.persistTicker.Stop()
	t.secondTicker.Stop()
	t.minuteTicker.Stop()
}

func (tf *PreparedTickFuncs) SecondTickRegister(f func(ctx Context)) {
	tf.secondTickFuncs = append(tf.secondTickFuncs, f)
}

func (tf *PreparedTickFuncs) MinuteTickRegister(f func(ctx Context)) {
	tf.minuteTickFuncs = append(tf.minuteTickFuncs, f)
}

func (tf *PreparedTickFuncs) OnLoadEventRegister(f func(ctx Context)) {
	tf.onLoadEventFuncs = append(tf.onLoadEventFuncs, f)
}

func (tf *PreparedTickFuncs) OnCreatedEventRegister(f func(ctx Context)) {
	tf.onCreatedEventFuncs = append(tf.onCreatedEventFuncs, f)
}

func (tf *PreparedTickFuncs) CustomEventRegister(e WorkerEventType, f eventFunc) {
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
