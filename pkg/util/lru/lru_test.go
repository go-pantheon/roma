package lru

import (
	"math/rand/v2"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// mockUserProto creates a mock user proto for testing
func mockUserProto(id int64) proto.Message {
	return &timestamppb.Timestamp{
		Seconds: id,
		Nanos:   int32(id),
	}
}

func TestLRU_Get(t *testing.T) {
	t.Parallel()

	lru := NewLRU(WithCapacity(2))

	ctime := time.Now()
	// Test getting non-existent key
	value, exists := lru.Get(1, ctime)
	assert.Nil(t, value)
	assert.False(t, exists)

	// Test getting existing key
	lru.Put(1, mockUserProto(1), ctime)
	value, exists = lru.Get(1, ctime)
	assert.NotNil(t, value)
	assert.True(t, exists)
	assert.Equal(t, int64(1), value.(*timestamppb.Timestamp).Seconds)
}

func TestLRU_Put(t *testing.T) {
	t.Parallel()

	lru := NewLRU(WithCapacity(2))
	ctime := time.Now()

	// Test putting new key
	lru.Put(1, mockUserProto(1), ctime)
	assert.Equal(t, 1, lru.Len())

	// Test updating existing key
	lru.Put(1, mockUserProto(2), ctime)
	assert.Equal(t, 1, lru.Len())
	value, _ := lru.Get(1, ctime)
	assert.Equal(t, int64(2), value.(*timestamppb.Timestamp).Seconds)

	// Test capacity limit
	lru.Put(2, mockUserProto(2), ctime)
	lru.Put(3, mockUserProto(3), ctime)
	assert.Equal(t, 2, lru.Len())
	_, exists := lru.Get(1, ctime)
	assert.False(t, exists)
}

func TestLRU_Clear(t *testing.T) {
	t.Parallel()

	lru := NewLRU(WithCapacity(2))
	ctime := time.Now()
	lru.Put(1, mockUserProto(1), ctime)
	lru.Put(2, mockUserProto(2), ctime)

	lru.Clear()
	assert.Equal(t, 0, lru.Len())

	// Verify cache is empty
	_, exists := lru.Get(1, ctime)
	assert.False(t, exists)
	_, exists = lru.Get(2, ctime)
	assert.False(t, exists)
}

func TestLRU_Concurrent(t *testing.T) {
	t.Parallel()

	lru := NewLRU(WithCapacity(100))
	done := make(chan bool)
	concurrent := 10

	ctime := time.Now()
	// Test concurrent writes
	for i := 0; i < concurrent; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				lru.Put(int64(id*100+j), mockUserProto(int64(id*100+j)), ctime)
			}
			done <- true
		}(i)
	}

	// Test concurrent reads
	for i := range concurrent {
		ctime := time.Now()

		go func(id int) {
			for j := range 100 {
				lru.Get(int64(id*100+j), ctime)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrent*2; i++ {
		<-done
	}

	// Verify final state
	assert.Equal(t, 100, lru.Len())
}

func TestLRU_EdgeCases(t *testing.T) {
	t.Parallel()

	ctime := time.Now()

	// Test with capacity 1
	lru := NewLRU(WithCapacity(1))
	lru.Put(1, mockUserProto(1), ctime)
	lru.Put(2, mockUserProto(2), ctime)
	assert.Equal(t, 1, lru.Len())
	_, exists := lru.Get(1, ctime)
	assert.False(t, exists)
	value, exists := lru.Get(2, ctime)
	assert.True(t, exists)
	assert.Equal(t, int64(2), value.(*timestamppb.Timestamp).Seconds)

	// Test with nil values
	lru.Clear()
	lru.Put(1, nil, ctime)
	value, exists = lru.Get(1, ctime)
	assert.Nil(t, value)
	assert.True(t, exists)
}

// Benchmark tests
func BenchmarkLRU_Get(b *testing.B) {
	size := 10_000
	ctime := time.Now()
	lru := NewLRU(WithCapacity(size))
	// Pre-fill the cache
	for i := 0; i < size; i++ {
		lru.Put(int64(i), mockUserProto(int64(i)), ctime)
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lru.Get(int64(rand.IntN(size)), ctime)
		}
	})
}

func BenchmarkLRU_Put(b *testing.B) {
	size := 10_000
	ctime := time.Now()
	lru := NewLRU(WithCapacity(size))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lru.Put(int64(rand.IntN(size)), mockUserProto(int64(rand.IntN(size))), ctime)
		}
	})
}

func BenchmarkLRU_Concurrent(b *testing.B) {
	size := 10_000
	ctime := time.Now()
	lru := NewLRU(WithCapacity(size))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lru.Put(int64(rand.IntN(size)), mockUserProto(int64(rand.IntN(size))), ctime)
			lru.Get(int64(rand.IntN(size)), ctime)
		}
	})
}

func BenchmarkLRU_Mixed(b *testing.B) {
	size := 1_000_000
	ctime := time.Now()
	lru := NewLRU(WithCapacity(size))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if rand.IntN(2) == 0 {
				lru.Put(int64(rand.IntN(size)), mockUserProto(int64(rand.IntN(size))), ctime)
			} else {
				lru.Get(int64(rand.IntN(size)), ctime)
			}
		}
	})
}
