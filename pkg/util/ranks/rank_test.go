package sorts

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
)

func TestSortTopN(t *testing.T) {
	tests := []struct {
		name     string
		arr      []int64
		top      int	
		less     func(i, j int64) bool
		expected []int64
	}{
		{
			name: "normal case ascending",
			arr:  []int64{5, 2, 8, 1, 9, 3},
			top:  3,
			less: func(i, j int64) bool { return i < j },
			expected: []int64{1, 2, 3},
		},
		{
			name: "normal case descending",
			arr:  []int64{5, 2, 8, 1, 9, 3},
			top:  3,
			less: func(i, j int64) bool { return i > j },
			expected: []int64{9, 8, 5},
		},
		{
			name:     "empty array",
			arr:      []int64{},
			top:      3,
			less:     func(i, j int64) bool { return i < j },
			expected: []int64{},
		},
		{
			name:     "top <= 0",
			arr:      []int64{1, 2, 3},
			top:      0,
			less:     func(i, j int64) bool { return i < j },
			expected: []int64{},
		},
		{
			name:     "top > len(arr)",
			arr:      []int64{3, 1, 2},
			top:      5,
			less:     func(i, j int64) bool { return i < j },
			expected: []int64{1, 2, 3},
		},
		{
			name:     "single element",
			arr:      []int64{1},
			top:      1,
			less:     func(i, j int64) bool { return i < j },
			expected: []int64{1},
		},
		{
			name:     "duplicate elements",
			arr:      []int64{3, 1, 3, 2, 3},
			top:      3,
			less:     func(i, j int64) bool { return i < j },
			expected: []int64{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortTopN(tt.arr, tt.top, tt.less)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortTopN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInsertIntoSortedArray(t *testing.T) {
	less := func(i, j int64) bool { return i < j }
	
	tests := []struct {
		name     string
		arr      []int64
		val      int64
		maxLen   int
		expected []int64
	}{
		{
			name:     "insert into empty array",
			arr:      []int64{},
			val:      5,
			maxLen:   3,
			expected: []int64{5},
		},
		{
			name:     "insert into non-full array",
			arr:      []int64{1, 3, 5},
			val:      2,
			maxLen:   5,
			expected: []int64{1, 2, 3, 5},
		},
		{
			name:     "insert into full array - should insert",
			arr:      []int64{1, 3, 5},
			val:      2,
			maxLen:   3,
			expected: []int64{1, 2, 3},
		},
		{
			name:     "insert into full array - should not insert",
			arr:      []int64{1, 3, 5},
			val:      6,
			maxLen:   3,
			expected: []int64{1, 3, 5},
		},
		{
			name:     "insert duplicate value",
			arr:      []int64{1, 2, 3},
			val:      2,
			maxLen:   4,
			expected: []int64{1, 2, 2, 3},
		},
		{
			name:     "insert at beginning",
			arr:      []int64{2, 3, 4},
			val:      1,
			maxLen:   3,
			expected: []int64{1, 2, 3},
		},
		{
			name:     "insert at end",
			arr:      []int64{1, 2},
			val:      3,
			maxLen:   3,
			expected: []int64{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a copy of input array to avoid modifying test data
			arr := make([]int64, len(tt.arr))
			copy(arr, tt.arr)
			
			result := InsertIntoSortedArray(arr, tt.val, tt.maxLen, less)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertIntoSortedArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSortTopN(b *testing.B) {
	sizes := []int{100, 1000, 10000}
	tops := []int{10, 100, 1000}

	for _, size := range sizes {
		for _, top := range tops {
			if top > size {
				continue
			}
			
			b.Run(fmt.Sprintf("size=%d,top=%d", size, top), func(b *testing.B) {
				arr := make([]int64, size)
				for i := range arr {
					arr[i] = rand.Int63()
				}
				less := func(i, j int64) bool { return i < j }
				
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					SortTopN(arr, top, less)
				}
			})
		}
	}
}

func BenchmarkInsertIntoSortedArray(b *testing.B) {
	sizes := []int{10, 100, 1000}
	
	for _, size := range sizes {
		b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
			arr := make([]int64, size)
			for i := range arr {
				arr[i] = int64(i)
			}
			less := func(i, j int64) bool { return i < j }
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// Create a copy of array for each iteration
				arrCopy := make([]int64, len(arr))
				copy(arrCopy, arr)
				InsertIntoSortedArray(arrCopy, rand.Int63(), size+1, less)
			}
		})
	}
} 