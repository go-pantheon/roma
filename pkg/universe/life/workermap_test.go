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
	t.Parallel()

	goroutines := 10
	iterations := 100

	m := NewWorkerMap()

	// Test Set and Get concurrently
	//nolint:paralleltest
	t.Run("Set and Get", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(goroutines)

		for i := range goroutines {
			go func(base int) {
				defer wg.Done()

				for j := range iterations {
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
	//nolint:paralleltest
	t.Run("Remove", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(goroutines)

		for i := range goroutines {
			go func(base int) {
				defer wg.Done()

				for j := range iterations {
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
	t.Parallel()

	count := 1000
	m := NewWorkerMap()

	var wg sync.WaitGroup

	// Concurrent Set
	//nolint:paralleltest
	t.Run("Concurrent Set", func(t *testing.T) {
		wg.Add(count)

		for i := range count {
			go func(i int64) {
				defer wg.Done()
				m.Set(i, &Worker{referer: fmt.Sprintf("referer-%d", i)})
			}(int64(i))
		}

		wg.Wait()
		assert.Equal(t, m.SimilarCount(), count)
	})

	// Concurrent Get
	//nolint:paralleltest
	t.Run("Concurrent Get", func(t *testing.T) {
		wg.Add(count)

		for i := range count {
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
	t.Parallel()

	goroutines := 10
	batchSize := 100

	m := NewWorkerMap()

	//nolint:paralleltest
	t.Run("Set", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(goroutines)

		for i := range goroutines {
			go func(base int) {
				defer wg.Done()

				for j := range batchSize {
					m.Set(int64(base*batchSize+j), &Worker{referer: fmt.Sprintf("referer-%d", int64(base*batchSize+j))})
				}
			}(i)
		}

		wg.Wait()
		assert.Equal(t, m.SimilarCount(), goroutines*batchSize)
	})

	//nolint:paralleltest
	t.Run("Get", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(goroutines)

		for i := range goroutines {
			go func(base int) {
				defer wg.Done()

				for j := range batchSize {
					val := m.Get(int64(base*batchSize + j))
					require.NotNil(t, val)
					assert.Equal(t, val.Referer(), fmt.Sprintf("referer-%d", int64(base*batchSize+j)))
				}
			}(i)
		}

		wg.Wait()
	})

	//nolint:paralleltest
	t.Run("Remove", func(t *testing.T) {
		var wg sync.WaitGroup

		wg.Add(goroutines)

		for i := range goroutines {
			go func(base int) {
				defer wg.Done()

				keys := make([]int64, batchSize)

				for j := range batchSize {
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

	for i := range 1000 {
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
			for range numCPU {
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
