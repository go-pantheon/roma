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
	ProductPreparedEvent(t WorkerEventType, args ...int64) error
	ProductFuncEvent(f EventFunc) error
	ConsumeEvent() <-chan EventFunc
}

type EventFunc func(wctx Context) (err error)
type eventFunc func(wctx Context, args ...int64) (err error)

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
