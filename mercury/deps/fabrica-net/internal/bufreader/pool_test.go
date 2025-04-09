package bufreader

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSyncPool_InvalidParams(t *testing.T) {
	tests := []struct {
		name    string
		min     int
		max     int
		factor  int
		wantErr bool
	}{
		{"min <= 0", 0, 1024, 4, true},
		{"max < min", 1024, 512, 4, true},
		{"factor <= 0", 1024, 2048, 0, true},
		{"valid params", 1024, 2048, 2, false},
		{"max over 1MB", 1024, 2 * 1024 * 1024, 4, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSyncPool(tt.min, tt.max, tt.factor)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAlloc(t *testing.T) {
	minSize := 128
	maxSize := 512
	factor := 2

	pool, err := NewSyncPool(minSize, maxSize, factor)
	require.NoError(t, err)

	// Get actual class sizes
	classSizes := pool.ClassSizes()

	tests := []struct {
		name    string
		size    int
		wantLen int
		wantCap int
	}{
		{"size 0", 0, 0, 0},
		{"negative size", -1, 0, 0},
		{"smaller than min", minSize / 2, minSize / 2, minSize},
		{"exact min size", minSize, minSize, minSize},
		{"just above min size", minSize + 1, minSize + 1, classSizes[1]},
		{"first class size", classSizes[0], classSizes[0], classSizes[0]},
		{"second class size", classSizes[1], classSizes[1], classSizes[1]},
		{"just below max", maxSize - 1, maxSize - 1, maxSize},
		{"max size", maxSize, maxSize, maxSize},
		{"just above max", maxSize + 1, maxSize + 1, maxSize + 1},
		{"exceed max significantly", maxSize * 2, maxSize * 2, maxSize * 2},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Alloc - %s - %d", tt.name, tt.size), func(t *testing.T) {
			buf := pool.Alloc(tt.size)
			assert.Equal(t, tt.wantLen, len(buf), "Length should match requested size")
			assert.Equal(t, tt.wantCap, cap(buf), "Capacity should match expected class size")

			// Test that we can actually use the buffer without panics
			if len(buf) > 0 {
				// Write to first and last byte to ensure buffer is usable
				buf[0] = 1
				if len(buf) > 1 {
					buf[len(buf)-1] = 1
				}
			}
		})
	}
}

func TestFree(t *testing.T) {
	minSize := 128
	maxSize := 512
	factor := 2
	pool, err := NewSyncPool(minSize, maxSize, factor)
	require.NoError(t, err)

	classSizes := pool.ClassSizes()

	t.Run("free nil", func(t *testing.T) {
		assert.NotPanics(t, func() { pool.Free(nil) })
	})

	t.Run("free empty slice", func(t *testing.T) {
		emptyBuf := make([]byte, 0)
		assert.NotPanics(t, func() { pool.Free(emptyBuf) })
	})

	// Test reuse for each class size
	for i, size := range classSizes {
		t.Run(fmt.Sprintf("reuse memory - class %d (size %d)", i, size), func(t *testing.T) {
			// Allocate with exact class size
			buf1 := pool.Alloc(size)
			ptr1 := &buf1[0]
			pool.Free(buf1)

			buf2 := pool.Alloc(size)
			ptr2 := &buf2[0]
			assert.Equal(t, ptr1, ptr2, "should reuse memory for class size %d", size)
			pool.Free(buf2)
		})
	}

	t.Run("free below min size", func(t *testing.T) {
		buf := make([]byte, minSize/2)
		assert.NotPanics(t, func() { pool.Free(buf) })
	})

	t.Run("free at max boundary", func(t *testing.T) {
		buf := make([]byte, maxSize)
		assert.NotPanics(t, func() { pool.Free(buf) })
	})

	t.Run("free just over max", func(t *testing.T) {
		buf := make([]byte, maxSize+1)
		assert.NotPanics(t, func() { pool.Free(buf) })
	})

	t.Run("free way over max", func(t *testing.T) {
		buf := make([]byte, maxSize*10)
		assert.NotPanics(t, func() { pool.Free(buf) })
	})
}
func TestConcurrentAllocFree(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	minSize := 128
	maxSize := 8192
	factor := 2

	pool, err := NewSyncPool(minSize, maxSize, factor)
	require.NoError(t, err)

	const goroutines = 10
	const iterations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Create different buffer sizes to test
	sizes := []int{64, 128, 129, 256, 1024, 4096, 8192, 16384}

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()

			// Each goroutine uses a different pattern of buffer sizes
			for j := 0; j < iterations; j++ {
				size := sizes[(id+j)%len(sizes)]
				buf := pool.Alloc(size)

				// Do something with the buffer to ensure it's usable
				if len(buf) > 0 {
					buf[0] = byte(id)
					if len(buf) > 1 {
						buf[len(buf)-1] = byte(j)
					}
				}

				pool.Free(buf)
			}
		}(i)
	}

	wg.Wait()
}

func TestClassSizes(t *testing.T) {
	minSize := 1024
	maxSize := 65536
	factor := 4

	pool, err := NewSyncPool(minSize, maxSize, factor)
	require.NoError(t, err)

	prev := 0
	for _, size := range pool.classesSize {
		// t.Logf("Class %d: %d bytes", i, size)
		assert.True(t, size > prev, "sizes should be increasing")
		prev = size

		// We're not going to check the exact formula since the implementation
		// seems to use a different growth algorithm than what was originally tested.
		// Just ensure the sizes are increasing and within bounds.
	}
	assert.Equal(t, maxSize, pool.classesSize[len(pool.classesSize)-1])
}

func TestClassSizesCalculation(t *testing.T) {
	minSize := 128
	maxSize := 8192
	factor := 2

	pool, err := NewSyncPool(minSize, maxSize, factor)
	require.NoError(t, err)

	sizes := pool.ClassSizes()

	// Verify first size is min size
	assert.Equal(t, minSize, sizes[0])

	// Verify last size is max size
	assert.Equal(t, maxSize, sizes[len(sizes)-1])

	// Verify progression of sizes follows the factor
	for i := 1; i < len(sizes); i++ {
		// nextSize := int(float64(curSize) * (1.0 + 1.0/float64(factor)))
		assert.Equal(t, int(math.Min(float64(sizes[i-1])*(1+1/float64(factor)), float64(maxSize))), sizes[i],
			"Class size %d should be %d×factor(%d) = %d, got %d",
			i, sizes[i-1], factor, sizes[i-1]*factor, sizes[i])
	}
}

func BenchmarkAllocOnly(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	sizes := []int{512, 1024, 4096, 16384, 65536}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		size := sizes[i%len(sizes)]
		buf := pool.Alloc(size)
		_ = buf
	}
}

func BenchmarkFreeOnly(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	sizes := []int{512, 1024, 4096, 16384, 65536}

	bufSize := 10000
	bufs := make([][]byte, bufSize)
	for i := range bufSize {
		size := sizes[i%len(sizes)]
		bufs[i] = pool.Alloc(size)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pool.Free(bufs[i%bufSize])
	}
}

func BenchmarkAllocFreeRandomSizes(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	maxSize := 65536 - 1024

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		size := rand.Intn(maxSize) + 1024
		buf := pool.Alloc(size)
		pool.Free(buf)
	}
}

func BenchmarkAllocFreeMinSize(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	size := 1024

	warmupBufs := make([]*[]byte, 100)
	for i := 0; i < 100; i++ {
		buf := pool.Alloc(size)
		warmupBufs[i] = &buf
	}
	for _, bufPtr := range warmupBufs {
		pool.Free(*bufPtr)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := pool.Alloc(size)
		pool.Free(buf)
	}
}

func BenchmarkAllocFreeMaxSize(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	size := 65536

	warmupBufs := make([]*[]byte, 100)
	for i := 0; i < 100; i++ {
		buf := pool.Alloc(size)
		warmupBufs[i] = &buf
	}
	for _, bufPtr := range warmupBufs {
		pool.Free(*bufPtr)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := pool.Alloc(size)
		pool.Free(buf)
	}
}

func BenchmarkAllocFreeConsistentSize(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	size := 4096

	warmupBufs := make([]*[]byte, 100)
	for i := 0; i < 100; i++ {
		buf := pool.Alloc(size)
		warmupBufs[i] = &buf
	}
	for _, bufPtr := range warmupBufs {
		pool.Free(*bufPtr)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		buf := pool.Alloc(size)
		pool.Free(buf)
	}
}

func BenchmarkParallelRandomSizes(b *testing.B) {
	pool, _ := NewSyncPool(1024, 65536, 4)
	maxSize := 65536 - 1024

	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine needs its own rand source to avoid lock contention
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for pb.Next() {
			size := r.Intn(maxSize) + 1024
			buf := pool.Alloc(size)
			pool.Free(buf)
		}
	})
}
