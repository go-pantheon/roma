package lru

import (
	"container/list"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"
)

// LRU implements a thread-safe Least Recently Used cache
type LRU struct {
	sync.RWMutex

	capacity int
	ttl      time.Duration
	ll       *list.List
	cache    map[int64]*list.Element
	onRemove func(key int64, value proto.Message)
}

type LRUOption func(*LRU)

func WithCapacity(capacity int) LRUOption {
	return func(l *LRU) {
		l.capacity = capacity
	}
}

func WithTTL(ttl time.Duration) LRUOption {
	return func(l *LRU) {
		l.ttl = ttl
	}
}

func WithOnRemove(onRemove func(key int64, value proto.Message)) LRUOption {
	return func(l *LRU) {
		l.onRemove = onRemove
	}
}

// entry represents a key-value pair in the cache
type entry struct {
	key   int64
	value proto.Message
	exp   time.Time
}

// NewLRU creates a new LRU cache with the given capacity
func NewLRU(opts ...LRUOption) *LRU {
	lru := &LRU{
		capacity: 5000, // default capacity
		ttl:      time.Minute,
		ll:       list.New(),
	}

	for _, opt := range opts {
		opt(lru)
	}

	if lru.capacity <= 0 {
		lru.capacity = 5000
	}

	if lru.ttl <= 0 {
		lru.ttl = time.Minute
	}

	lru.cache = make(map[int64]*list.Element, lru.capacity)

	return lru
}

// Get retrieves a value from the cache
func (l *LRU) Get(key int64, ctime time.Time) (proto.Message, bool) {
	l.Lock()
	defer l.Unlock()

	if ele, hit := l.cache[key]; hit {
		if ele.Value.(*entry).exp.Before(ctime) {
			l.ll.Remove(ele)
			delete(l.cache, key)

			return nil, false
		}

		l.ll.MoveToFront(ele)

		return ele.Value.(*entry).value, true
	}

	return nil, false
}

// Put adds or updates a value in the cache
func (l *LRU) Put(key int64, value proto.Message, ctime time.Time) {
	l.Lock()
	defer l.Unlock()

	if ele, ok := l.cache[key]; ok {
		l.ll.MoveToFront(ele)

		oldValue := ele.Value.(*entry).value

		ele.Value.(*entry).value = value
		ele.Value.(*entry).exp = ctime.Add(l.ttl)

		if l.onRemove != nil {
			l.onRemove(key, oldValue)
		}

		return
	}

	l.cache[key] = l.ll.PushFront(&entry{key: key, value: value, exp: ctime.Add(l.ttl)})

	if l.ll.Len() > l.capacity {
		l.removeOldest()
	}
}

// removeOldest removes the least recently used item from the cache
func (l *LRU) removeOldest() {
	ele := l.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)

		l.ll.Remove(ele)
		delete(l.cache, kv.key)

		if l.onRemove != nil {
			l.onRemove(kv.key, kv.value)
		}
	}
}

func (l *LRU) Remove(key int64) {
	l.Lock()
	defer l.Unlock()

	if ele, ok := l.cache[key]; ok {
		l.ll.Remove(ele)
		delete(l.cache, key)
	}
}

// Len returns the current number of items in the cache
func (l *LRU) Len() int {
	l.RLock()
	defer l.RUnlock()

	return l.ll.Len()
}

// Clear removes all items from the cache
func (l *LRU) Clear() {
	l.Lock()
	defer l.Unlock()

	l.ll = list.New()
	l.cache = make(map[int64]*list.Element)
}
