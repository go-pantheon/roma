package life

import (
	"sync"
	"sync/atomic"
)

const shardCount = uint64(512)

// WorkerMap is a thread-safe map implementation that uses sharding to reduce lock contention.
// It divides the map into multiple shards, each protected by its own RWMutex.
// This allows for better concurrent access compared to a single map protected by a single lock.
type WorkerMap struct {
	shards []*mapShared
}

type options struct {
	shardCount uint64
}

type Option func(*options)

func WithShardCount(shardCount uint64) Option {
	return func(o *options) {
		o.shardCount = shardCount
	}
}

func NewWorkerMap(opts ...Option) *WorkerMap {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

	if opt.shardCount == 0 {
		opt.shardCount = shardCount
	}

	m := &WorkerMap{
		shards: make([]*mapShared, opt.shardCount),
	}
	for i := range opt.shardCount {
		m.shards[i] = newMapShared()
	}

	return m
}

func (m WorkerMap) Set(key int64, value *Worker) (old *Worker) {
	shard := m.getShard(key)

	return shard.loadAndStore(key, value)
}

func (m WorkerMap) Get(key int64) *Worker {
	return m.getShard(key).get(key)
}

func (m WorkerMap) SimilarCount() int {
	count := 0
	for i := range shardCount {
		count += int(m.shards[i].similarCount())
	}

	return count
}

func (m WorkerMap) Remove(key int64) {
	m.getShard(key).delete(key)
}

// Tuple used by the Iter functions to wrap two variables together over a channel,
type Tuple struct {
	Key int64
	Val *Worker
}

// Iter returns a buffered iterator which could be used in a for range loop.
func (m WorkerMap) Iter() <-chan Tuple {
	total := 0

	chans := snapshot(m)
	for _, c := range chans {
		total += cap(c)
	}

	ch := make(chan Tuple, total)

	go fanIn(chans, ch)

	return ch
}

// fanIn reads elements from channels `chans` into channel `out`
func fanIn(chans []chan Tuple, out chan Tuple) {
	wg := sync.WaitGroup{}

	wg.Add(len(chans))

	for _, ch := range chans {
		go func(ch chan Tuple) {
			for t := range ch {
				out <- t
			}

			wg.Done()
		}(ch)
	}

	wg.Wait()
	close(out)
}

// snapshot returns an array of channels that contains elements in each shard,
// which likely takes a snapshot of `m`.
// It returns once the size of each buffered channel is determined,
// before all the channels are populated using goroutines.
func snapshot(m WorkerMap) (chans []chan Tuple) {
	chans = make([]chan Tuple, shardCount)
	for index, shard := range m.shards {
		chans[index] = make(chan Tuple, int(shard.similarCount()))
		shard.workers.Range(func(key, value any) bool {
			chans[index] <- Tuple{key.(int64), value.(*Worker)}
			return true
		})
		close(chans[index])
	}

	return chans
}

func (m WorkerMap) getShard(key int64) *mapShared {
	return m.shards[getShardIndex(key)]
}

func getShardIndex(key int64) uint64 {
	return wyhash(key) & (shardCount - 1)
}

// wyhash generates a 64-bit hash for the given 64-bit key using wyhash algorithm.
func wyhash(key int64) uint64 {
	x := uint64(key)
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33

	return x
}

type mapShared struct {
	workers *sync.Map
	count   atomic.Int64
}

func newMapShared() *mapShared {
	return &mapShared{
		workers: &sync.Map{},
	}
}

func (s *mapShared) get(key int64) *Worker {
	if val, ok := s.workers.Load(key); ok {
		return val.(*Worker)
	}

	return nil
}

func (s *mapShared) loadAndStore(key int64, value *Worker) (old *Worker) {
	if val, ok := s.workers.LoadOrStore(key, value); ok {
		old = val.(*Worker)

		s.workers.Store(key, value)
	} else {
		s.count.Add(1)
	}

	return old
}

func (s *mapShared) delete(key int64) {
	if _, ok := s.workers.LoadAndDelete(key); ok {
		s.count.Add(-1)
	}
}

func (s *mapShared) similarCount() int64 {
	return s.count.Load()
}
