package life

import (
	"time"

	"github.com/go-pantheon/roma/pkg/universe/constants"
)

type Tickers struct {
	*BuiltinEventFuncs

	renewalTicker *time.Ticker
	persistTicker *time.Ticker
	secondTicker  *time.Ticker
	minuteTicker  *time.Ticker
}

func newTickers() *Tickers {
	t := &Tickers{
		renewalTicker: time.NewTicker(constants.WorkerRenewalTickDuration),
		secondTicker:  time.NewTicker(constants.WorkerSecondTickDuration),
		minuteTicker:  time.NewTicker(constants.WorkerMinuteTickDuration),
		persistTicker: time.NewTicker(constants.WorkerPersistTickDuration),
	}

	return t
}

func (t *Tickers) stop() {
	t.renewalTicker.Stop()
	t.persistTicker.Stop()
	t.secondTicker.Stop()
	t.minuteTicker.Stop()
}
