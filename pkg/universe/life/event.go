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

type WithArg func(arg *EventArg)

func WithI64(i64 int64) WithArg {
	return func(arg *EventArg) {
		arg.i64 = i64
	}
}

func WithStr(str string) WithArg {
	return func(arg *EventArg) {
		arg.str = str
	}
}

func WithF64(f64 float64) WithArg {
	return func(arg *EventArg) {
		arg.f64 = f64
	}
}

func WithI64s(i64s ...int64) WithArg {
	return func(arg *EventArg) {
		arg.i64s = i64s
	}
}

func WithStrs(strs ...string) WithArg {
	return func(arg *EventArg) {
		arg.strs = strs
	}
}

func WithF64s(f64s ...float64) WithArg {
	return func(arg *EventArg) {
		arg.f64s = f64s
	}
}

type EventArg struct {
	i64  int64
	str  string
	f64  float64
	i64s []int64
	strs []string
	f64s []float64
}

func (e *EventArg) Reset() {
	e.i64 = 0
	e.str = ""
	e.f64 = 0
	e.i64s = nil
	e.strs = nil
	e.f64s = nil
}

func (e *EventArg) I64() int64 {
	return e.i64
}

func (e *EventArg) Str() string {
	return e.str
}

func (e *EventArg) F64() float64 {
	return e.f64
}

func (e *EventArg) I64s() []int64 {
	return e.i64s
}

func (e *EventArg) Strs() []string {
	return e.strs
}

func (e *EventArg) F64s() []float64 {
	return e.f64s
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
