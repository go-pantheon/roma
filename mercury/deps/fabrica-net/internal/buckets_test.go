package internal

import (
	"sync"
	"testing"

	net "github.com/go-pantheon/fabrica-net"
	"github.com/go-pantheon/fabrica-net/conf"
	"github.com/stretchr/testify/assert"
)

func TestNewBuckets(t *testing.T) {
	tests := []struct {
		name       string
		bucketSize int
		wantSize   int
	}{
		{
			name:       "normal size",
			bucketSize: 16,
			wantSize:   16,
		},
		{
			name:       "small size",
			bucketSize: 1,
			wantSize:   1,
		},
		{
			name:       "large size",
			bucketSize: 1024,
			wantSize:   1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &conf.Bucket{
				BucketSize: tt.bucketSize,
			}
			bs := NewBuckets(c)
			assert.Equal(t, tt.wantSize, len(bs.buckets))
			assert.Equal(t, uint64(tt.bucketSize), shardCount)
		})
	}
}

func TestBuckets_BasicOperations(t *testing.T) {
	c := &conf.Bucket{
		BucketSize: 16,
	}
	bs := NewBuckets(c)

	// Test Put and Get
	w1 := newTestWorker(1, 100)
	w2 := newTestWorker(2, 200)

	// Test Put
	old := bs.Put(w1)
	assert.Nil(t, old)
	old = bs.Put(w2)
	assert.Nil(t, old)

	// Test Get
	got := bs.Worker(w1.WID())
	assert.Equal(t, w1, got)
	got = bs.Worker(w2.WID())
	assert.Equal(t, w2, got)

	// Test Del
	bs.Del(w1)
	got = bs.Worker(w1.WID())
	assert.Nil(t, got)
}

func TestBuckets_UIDOperations(t *testing.T) {
	c := &conf.Bucket{
		BucketSize: 16,
	}
	bs := NewBuckets(c)

	w1 := newTestWorker(1, 100)
	w2 := newTestWorker(2, 200)

	bs.Put(w1)
	bs.Put(w2)

	// Test GetByUID
	got := bs.GetByUID(w1.UID())
	assert.Equal(t, w1, got)
	got = bs.GetByUID(w2.UID())
	assert.Equal(t, w2, got)
	got = bs.GetByUID(999) // non-existent UID
	assert.Nil(t, got)

	// Test GetByUIDs
	workers := bs.GetByUIDs([]int64{w1.UID(), w2.UID(), 999})
	assert.Equal(t, 2, len(workers))
	assert.Contains(t, workers, w1.UID())
	assert.Contains(t, workers, w2.UID())
}

func TestBuckets_ConcurrentOperations(t *testing.T) {
	c := &conf.Bucket{
		BucketSize: 16,
	}
	bs := NewBuckets(c)

	// Create multiple workers
	workers := make([]*Worker, 100)
	for i := 0; i < 100; i++ {
		workers[i] = newTestWorker(uint64(i), int64(i))
	}

	// Concurrent Put operations
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(w *Worker) {
			bs.Put(w)
			done <- true
		}(workers[i])
	}

	// Wait for all goroutines to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify all workers are present
	for _, w := range workers {
		got := bs.Worker(w.WID())
		assert.Equal(t, w, got)
	}
}

func TestBuckets_Walk(t *testing.T) {
	c := &conf.Bucket{
		BucketSize: 16,
	}
	bs := NewBuckets(c)

	// Add some workers
	workers := make([]*Worker, 5)
	for i := 0; i < 5; i++ {
		workers[i] = newTestWorker(uint64(i), int64(i))
		bs.Put(workers[i])
	}

	// Test Walk
	count := 0
	bs.Walk(func(w *Worker) bool {
		count++
		return true
	})
	assert.Equal(t, 5, count)

	// Test Walk with early termination
	count = 0
	bs.Walk(func(w *Worker) bool {
		count++
		return count < 3
	})
	assert.Equal(t, 3, count)
}

func BenchmarkBuckets_Operations(b *testing.B) {
	c := &conf.Bucket{
		BucketSize: 16,
	}
	bs := NewBuckets(c)

	// Prepare test data
	workers := make([]*Worker, 1000)
	for i := 1; i <= 1001; i++ {
		workers[i] = newTestWorker(uint64(i), int64(i))
		bs.Put(workers[i])
	}

	b.Run("Put", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			w := newTestWorker(uint64(i), int64(i))
			bs.Put(w)
		}
	})

	b.Run("Get", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bs.Worker(workers[i%1000].WID())
		}
	})

	b.Run("GetByUID", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bs.GetByUID(workers[i%1000].UID())
		}
	})

	b.Run("Walk", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			bs.Walk(func(w *Worker) bool {
				return true
			})
		}
	})
}

func BenchmarkBuckets_Concurrent(b *testing.B) {
	conf.Init()
	c := conf.Conf.Bucket
	bs := NewBuckets(c)

	// Prepare test data
	workers := make([]*Worker, 1000)
	for i := 0; i < 1000; i++ {
		workers[i] = newTestWorker(uint64(i), int64(i))
		bs.Put(workers[i])
	}

	b.Run("ConcurrentPut", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w := newTestWorker(uint64(b.N), int64(b.N))
				bs.Put(w)
			}
		})
	})

	b.Run("ConcurrentGet", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs.Worker(workers[b.N%1000].WID())
			}
		})
	})

	b.Run("ConcurrentGetByUID", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs.GetByUID(workers[b.N%1000].UID())
			}
		})
	})

	b.Run("ConcurrentWalk", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				bs.Walk(func(w *Worker) bool { return true })
			}
		})
	})
}

func BenchmarkMap_Comparison(b *testing.B) {
	// Prepare test data
	workers := make([]*Worker, 0, 20_000)
	for i := 1; i <= 20_000; i++ {
		workers = append(workers, newTestWorker(uint64(i), int64(i)))
	}

	// Standard map with mutex
	type mutexMap struct {
		sync.RWMutex
		m    map[uint64]*Worker
		uids *sync.Map
	}

	// Initialize map
	m := &mutexMap{
		m:    make(map[uint64]*Worker, 512*128),
		uids: &sync.Map{},
	}

	for i := 0; i < len(workers); i++ {
		m.m[uint64(i)] = workers[i]
		m.uids.Store(int64(i), uint64(i))
	}

	b.Run("MutexMapPut", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w := newTestWorker(uint64(b.N), int64(b.N))
				m.Lock()
				m.m[w.WID()] = w
				m.uids.Store(w.UID(), w.WID())
				m.Unlock()
			}
		})
	})

	b.Run("MutexMapGet", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.RLock()
				_ = m.m[workers[b.N%len(workers)].WID()]
				m.uids.Load(int64(b.N))
				m.RUnlock()
			}
		})
	})

	b.Run("MutexMapGetByUID", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.RLock()
				wid, _ := m.uids.Load(b.N % len(workers))
				if wid != nil {
					_ = m.m[wid.(uint64)]
				}

				m.RUnlock()
			}
		})
	})

	b.Run("MutexMapWalk", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				m.RLock()
				for _, w := range m.m {
					_ = w
				}
				m.RUnlock()
			}
		})
	})

	// sync.Map comparison
	type syncMap struct {
		m    *sync.Map
		uids *sync.Map
	}

	sm := &syncMap{
		m:    &sync.Map{},
		uids: &sync.Map{},
	}

	for i := 0; i < len(workers); i++ {
		sm.m.Store(workers[i].WID(), workers[i])
		sm.uids.Store(workers[i].UID(), workers[i].WID())
	}

	b.Run("SyncMapPut", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w := newTestWorker(uint64(b.N), int64(b.N))
				sm.m.Store(w.WID(), w)
				sm.uids.Store(w.UID(), w.WID())
			}
		})
	})

	b.Run("SyncMapGet", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				any, _ := sm.m.Load(workers[b.N%len(workers)].WID())
				_ = any.(*Worker)
				sm.uids.Load(int64(b.N))
			}
		})
	})

	b.Run("SyncMapGetByUID", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wid, _ := sm.uids.Load(b.N % len(workers))
				if wid != nil {
					any, _ := sm.m.Load(wid.(uint64))
					_ = any.(*Worker)
				}
			}
		})
	})

	b.Run("SyncMapWalk", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				sm.m.Range(func(key, value any) bool {
					_ = value.(*Worker)
					return true
				})
			}
		})
	})

	c := &conf.Bucket{
		BucketSize: 128,
	}

	buckets := NewBuckets(c)

	for i := 0; i < len(workers); i++ {
		buckets.Put(workers[i])
	}

	b.Run("BucketsPut", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				buckets.Put(newTestWorker(uint64(b.N), int64(b.N)))
			}
		})
	})

	b.Run("BucketsGet", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				buckets.Worker(workers[b.N%len(workers)].WID())
			}
		})
	})

	b.Run("BucketsGetByUID", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				buckets.GetByUID(workers[b.N%len(workers)].UID())
			}
		})
	})

	b.Run("BucketsWalk", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				buckets.Walk(func(w *Worker) bool { return true })
			}
		})
	})
}

func TestBuckets_ConcurrencySafety(t *testing.T) {
	c := &conf.Bucket{
		BucketSize: 16,
	}
	bs := NewBuckets(c)

	// test parameters
	const (
		numWorkers    = 1000 // total worker count
		numGoroutines = 100  // concurrent goroutine count
		numOperations = 1000 // operations per goroutine
	)

	// create test data
	workers := make([]*Worker, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = newTestWorker(uint64(i), int64(i))
	}

	// concurrent Put test
	t.Run("ConcurrentPut", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// start multiple goroutines to perform concurrent Put
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					workerID := (routineID*numOperations + j) % numWorkers
					bs.Put(workers[workerID])
				}
			}(i)
		}
		wg.Wait()

		// verify all workers are correctly stored
		for _, w := range workers {
			got := bs.Worker(w.WID())
			assert.Equal(t, w, got, "Worker not found after concurrent Put")
		}
	})

	// concurrent Get test
	t.Run("ConcurrentGet", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// start multiple goroutines to perform concurrent Get
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					workerID := (routineID*numOperations + j) % numWorkers
					got := bs.Worker(workers[workerID].WID())
					assert.Equal(t, workers[workerID], got, "Incorrect worker retrieved during concurrent Get")
				}
			}(i)
		}
		wg.Wait()
	})

	// concurrent GetByUID test
	t.Run("ConcurrentGetByUID", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// start multiple goroutines to perform concurrent GetByUID
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					workerID := (routineID*numOperations + j) % numWorkers
					got := bs.GetByUID(workers[workerID].UID())
					assert.Equal(t, workers[workerID], got, "Incorrect worker retrieved during concurrent GetByUID")
				}
			}(i)
		}
		wg.Wait()
	})

	// concurrent Put and Get mixed test
	t.Run("ConcurrentPutAndGet", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines * 2) // half for Put, half for Get

		// start Put goroutines
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					workerID := (routineID*numOperations + j) % numWorkers
					bs.Put(workers[workerID])
				}
			}(i)
		}

		// start Get goroutines
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					workerID := (routineID*numOperations + j) % numWorkers
					got := bs.Worker(workers[workerID].WID())
					if got != nil {
						assert.Equal(t, workers[workerID], got, "Incorrect worker retrieved during concurrent Put and Get")
					}
				}
			}(i)
		}
		wg.Wait()
	})

	// concurrent Del test
	t.Run("ConcurrentDel", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// start multiple goroutines to perform concurrent Del
		for i := 0; i < numGoroutines; i++ {
			go func(routineID int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					workerID := (routineID*numOperations + j) % numWorkers
					bs.Del(workers[workerID])
				}
			}(i)
		}
		wg.Wait()

		// verify all workers are correctly deleted
		for _, w := range workers {
			got := bs.Worker(w.WID())
			assert.Nil(t, got, "Worker still exists after concurrent Del")
		}
	})

	// concurrent Walk test
	t.Run("ConcurrentWalk", func(t *testing.T) {
		// add some workers again
		for _, w := range workers[:numWorkers/2] {
			bs.Put(w)
		}

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		// start multiple goroutines to perform concurrent Walk
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				count := 0
				bs.Walk(func(w *Worker) bool {
					count++
					return true
				})
				assert.Equal(t, numWorkers/2, count, "Incorrect count during concurrent Walk")
			}()
		}
		wg.Wait()
	})
}

func newTestWorker(id uint64, uid int64) *Worker {
	return &Worker{
		id:      id,
		session: net.NewSession(uid, 0, 0, nil, nil, false, "", 0),
	}
}
