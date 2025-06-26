package life

import "github.com/go-pantheon/fabrica-util/errors"

type BuiltinEventFuncs struct {
	secondTickFuncs     []func(ctx Context) error
	minuteTickFuncs     []func(ctx Context) error
	onCreatedEventFuncs []func(ctx Context) error
	onLoadEventFuncs    []func(ctx Context) error
}

func newBuiltinEventFuncs() *BuiltinEventFuncs {
	return &BuiltinEventFuncs{}
}

func (tf *BuiltinEventFuncs) RegisterSecondTick(f func(ctx Context) error) {
	tf.secondTickFuncs = append(tf.secondTickFuncs, f)
}

func (tf *BuiltinEventFuncs) RegisterMinuteTick(f func(ctx Context) error) {
	tf.minuteTickFuncs = append(tf.minuteTickFuncs, f)
}

func (tf *BuiltinEventFuncs) RegisterOnLoadEvent(f func(ctx Context) error) {
	tf.onLoadEventFuncs = append(tf.onLoadEventFuncs, f)
}

func (tf *BuiltinEventFuncs) RegisterOnCreatedEvent(f func(ctx Context) error) {
	tf.onCreatedEventFuncs = append(tf.onCreatedEventFuncs, f)
}

func (tf *BuiltinEventFuncs) RegisterCustomEvent(e WorkerEventType, f eventFunc) {
	customEventFuncMap.put(e, f)
}

func (tf *BuiltinEventFuncs) secondTick(wctx Context) (err error) {
	for _, f := range tf.secondTickFuncs {
		if tickErr := f(wctx); tickErr != nil {
			err = errors.Join(err, tickErr)
		}
	}

	return err
}

func (tf *BuiltinEventFuncs) minuteTick(wctx Context) (err error) {
	for _, f := range tf.minuteTickFuncs {
		if tickErr := f(wctx); tickErr != nil {
			err = errors.Join(err, tickErr)
		}
	}

	return err
}
