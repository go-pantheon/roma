package life

import (
	"sync"
)

var (
	preparedEventFuncMap *eventFuncMap
)

func init() {
	preparedEventFuncMap = newPreparedEventMap()
}

// WorkerEventType worker event type. Custom event types should be less than 10_000
type WorkerEventType int32

// Start from 10_000. The previous numbers are reserved for custom event types: v1.WorkerEventType
const (
	EventTypeSecondTick = WorkerEventType(iota + 10_000)
	EventTypeMinuteTick
)

type EventManageable interface {
	ProductPreparedEvent(t WorkerEventType, args ...WithArg) error
	ProductFuncEvent(f EventFunc) error
	ConsumeEvent() <-chan EventFunc
}

type EventFunc func(wctx Context) (err error)
type eventFunc func(wctx Context, arg *EventArg) (err error)

type eventFuncMap struct {
	sync.RWMutex

	funcs map[WorkerEventType][]eventFunc
}

func newPreparedEventMap() *eventFuncMap {
	m := &eventFuncMap{
		funcs: make(map[WorkerEventType][]eventFunc, 32),
	}
	return m
}

func (m *eventFuncMap) put(t WorkerEventType, f eventFunc) {
	m.Lock()
	defer m.Unlock()

	m.funcs[t] = append(m.funcs[t], f)
}

func (m *eventFuncMap) get(t WorkerEventType) ([]eventFunc, bool) {
	m.RLock()
	defer m.RUnlock()

	fs, ok := m.funcs[t]
	return fs, ok
}

type EventArgValue interface {
	int64 | string | float64 | []int64 | []string | []float64
}

type WithArg func(arg *EventArg)

func With[T EventArgValue](key any, value T) WithArg {
	return func(arg *EventArg) {
		arg.set(key, value)
	}
}

type EventArg struct {
	data map[any]any
}

func (e *EventArg) set(key any, value any) {
	if e.data == nil {
		e.data = make(map[any]any)
	}

	e.data[key] = value
}

func GetArg[T EventArgValue](arg *EventArg, key any) T {
	if arg.data == nil {
		var zero T
		return zero
	}

	v, ok := arg.data[key].(T)
	if !ok {
		var zero T
		return zero
	}

	return v
}

func (e *EventArg) Reset() {
	clear(e.data)
}

var eventArgPool = sync.Pool{
	New: func() any {
		return &EventArg{}
	},
}

func GetEventArg() *EventArg {
	return eventArgPool.Get().(*EventArg)
}

func PutEventArg(e *EventArg) {
	e.Reset()
	eventArgPool.Put(e)
}
