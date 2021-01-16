package life

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkMap_Basic(t *testing.T) {
	m := NewWorkerMap()
	const goroutines = 10
	const iterations = 100

	// Test Set and Get concurrently
	t.Run("Set and Get", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(base int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					key := int64(base*iterations + j)
					referer := fmt.Sprintf("referer-%d", key)
					m.Set(key, &Worker{referer: referer})
					val := m.Get(key)
					require.NotNil(t, val)
					assert.Equal(t, val.Referer(), referer)
				}
			}(i)
		}
		wg.Wait()
	})

	// Test Remove concurrently
	t.Run("Remove", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(base int) {
				defer wg.Done()
				for j := 0; j < iterations; j++ {
					key := int64(base*iterations + j)
					m.Remove(key)
					assert.Nil(t, m.Get(key))
				}
			}(i)
		}
		wg.Wait()
	})
}

func TestWorkMap_Concurrent(t *testing.T) {
	m := NewWorkerMap()
	count := 1000
	var wg sync.WaitGroup

	// Concurrent Set
	t.Run("Concurrent Set", func(t *testing.T) {
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(i int64) {
				defer wg.Done()
				m.Set(i, &Worker{referer: fmt.Sprintf("referer-%d", i)})
			}(int64(i))
		}
		wg.Wait()

		assert.Equal(t, m.Count(), count)
	})

	// Concurrent Get
	t.Run("Concurrent Get", func(t *testing.T) {
		wg.Add(count)
		for i := 0; i < count; i++ {
			go func(i int64) {
				defer wg.Done()
				val := m.Get(i)
				require.NotNil(t, val)
				assert.Equal(t, val.Referer(), fmt.Sprintf("referer-%d", i))
			}(int64(i))
		}
		wg.Wait()
	})
}

func TestWorkMap_BatchOperations(t *testing.T) {
	m := NewWorkerMap()
	const goroutines = 10
	const batchSize = 100

	t.Run("Set", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(base int) {
				defer wg.Done()
				for j := 0; j < batchSize; j++ {
					m.Set(int64(base*batchSize+j), &Worker{referer: fmt.Sprintf("referer-%d", int64(base*batchSize+j))})
				}
			}(i)
		}
		wg.Wait()
		assert.Equal(t, m.Count(), goroutines*batchSize)
	})

	t.Run("Get", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(base int) {
				defer wg.Done()
				for j := 0; j < batchSize; j++ {
					val := m.Get(int64(base*batchSize + j))
					require.NotNil(t, val)
					assert.Equal(t, val.Referer(), fmt.Sprintf("referer-%d", int64(base*batchSize+j)))
				}
			}(i)
		}
		wg.Wait()
	})

	t.Run("Remove", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for i := 0; i < goroutines; i++ {
			go func(base int) {
				defer wg.Done()
				keys := make([]int64, batchSize)
				for j := 0; j < batchSize; j++ {
					keys[j] = int64(base*batchSize + j)
				}
				for _, k := range keys {
					m.Remove(k)
					assert.Nil(t, m.Get(k))
				}
			}(i)
		}
		wg.Wait()
	})
}

// Benchmarks

func BenchmarkWorkMap_Set(b *testing.B) {
	m := NewWorkerMap()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set(rand.Int63(), &Worker{referer: fmt.Sprintf("referer-%d", rand.Int63())})
		}
	})
}

func BenchmarkWorkMap_Get(b *testing.B) {
	m := NewWorkerMap()
	for i := 0; i < 1000; i++ {
		m.Set(int64(i), &Worker{referer: fmt.Sprintf("referer-%d", i)})
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get(rand.Int63n(1000))
		}
	})
}

// Heavy load test
func BenchmarkWorkMap_HeavyLoad(b *testing.B) {
	m := NewWorkerMap()
	numCPU := runtime.NumCPU()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for i := 0; i < numCPU; i++ {
				key := rand.Int63()
				switch rand.Intn(3) {
				case 0:
					m.Set(key, &Worker{referer: fmt.Sprintf("referer-%d", key)})
				case 1:
					m.Get(key)
				case 2:
					m.Remove(key)
				}
			}
		}
	})
}

// Test different shard sizes
func BenchmarkWorkMap_DifferentShardSizes(b *testing.B) {
	sizes := []uint64{128, 512, 2048, 8192}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("ShardSize_%d", size), func(b *testing.B) {
			m := NewWorkerMap(WithShardCount(size))
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					key := rand.Int63()
					m.Set(key, &Worker{referer: fmt.Sprintf("referer-%d", key)})
					m.Get(key)
				}
			})
		})
	}
}
