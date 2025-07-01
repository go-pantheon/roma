package u64

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []uint64
		value    uint64
		expected bool
	}{
		{
			name:     "empty array",
			array:    []uint64{},
			value:    1,
			expected: false,
		},
		{
			name:     "value exists",
			array:    []uint64{1, 2, 3, 4, 5},
			value:    3,
			expected: true,
		},
		{
			name:     "value does not exist",
			array:    []uint64{1, 2, 3, 4, 5},
			value:    6,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := Contains(tt.array, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsElementRepeat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []uint64
		except   uint64
		expected bool
	}{
		{
			name:     "empty array",
			array:    []uint64{},
			except:   0,
			expected: false,
		},
		{
			name:     "no repeat",
			array:    []uint64{1, 2, 3, 4},
			except:   0,
			expected: false,
		},
		{
			name:     "has repeat",
			array:    []uint64{1, 2, 2, 3},
			except:   0,
			expected: true,
		},
		{
			name:     "has repeat but excepted",
			array:    []uint64{1, 2, 2, 3},
			except:   2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := IsElementRepeat(tt.array, tt.except)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIndex(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []uint64
		value    uint64
		expected int
	}{
		{
			name:     "empty array",
			array:    []uint64{},
			value:    1,
			expected: -1,
		},
		{
			name:     "value exists",
			array:    []uint64{1, 2, 3, 4, 5},
			value:    3,
			expected: 2,
		},
		{
			name:     "value does not exist",
			array:    []uint64{1, 2, 3, 4, 5},
			value:    6,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := Index(tt.array, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []uint64
		count    uint64
		expected int
	}{
		{
			name:     "empty array",
			array:    []uint64{},
			count:    5,
			expected: 0,
		},
		{
			name:     "count less than array length",
			array:    []uint64{1, 2, 3, 4, 5},
			count:    3,
			expected: 3,
		},
		{
			name:     "count greater than array length",
			array:    []uint64{1, 2, 3},
			count:    5,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := Rand(tt.array, tt.count)
			assert.Equal(t, tt.expected, len(result))
		})
	}
}

func TestCycle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		i           uint64
		array       []uint64
		expected    uint64
		expectError bool
	}{
		{
			name:        "empty array",
			i:           1,
			array:       []uint64{},
			expectError: true,
		},
		{
			name:        "i = 0",
			i:           0,
			array:       []uint64{1, 2, 3},
			expected:    1,
			expectError: false,
		},
		{
			name:        "normal cycle",
			i:           4,
			array:       []uint64{1, 2, 3},
			expected:    2,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Cycle(tt.i, tt.array)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestToKVMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       []uint64
		expected    map[uint64]uint64
		expectError bool
	}{
		{
			name:        "odd length",
			input:       []uint64{1, 2, 3},
			expectError: true,
		},
		{
			name:     "empty array",
			input:    []uint64{},
			expected: map[uint64]uint64{},
		},
		{
			name:  "normal case",
			input: []uint64{1, 10, 2, 20},
			expected: map[uint64]uint64{
				1: 10,
				2: 20,
			},
		},
		{
			name:  "duplicate keys",
			input: []uint64{1, 10, 1, 20},
			expected: map[uint64]uint64{
				1: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := ToKVMap(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestF64ArraysToI64Arrays(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    [][]float64
		expected [][]uint64
	}{
		{
			name:     "empty array",
			input:    [][]float64{},
			expected: [][]uint64{},
		},
		{
			name: "normal case",
			input: [][]float64{
				{1.0, 2.0, 3.0},
				{4.0, 5.0, 6.0},
			},
			expected: [][]uint64{
				{1, 2, 3},
				{4, 5, 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := F64ArraysToI64Arrays(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFirst(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		array       []uint64
		expected    uint64
		expectError bool
	}{
		{
			name:        "nil array",
			array:       nil,
			expectError: true,
		},
		{
			name:        "empty array",
			array:       []uint64{},
			expectError: true,
		},
		{
			name:        "normal case",
			array:       []uint64{1, 2, 3},
			expected:    1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := First(tt.array)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		array       []uint64
		index       int
		expected    uint64
		expectError bool
	}{
		{
			name:        "nil array",
			array:       nil,
			index:       0,
			expectError: true,
		},
		{
			name:        "empty array",
			array:       []uint64{},
			index:       0,
			expectError: true,
		},
		{
			name:        "index out of range",
			array:       []uint64{1, 2, 3},
			index:       3,
			expectError: true,
		},
		{
			name:        "normal case",
			array:       []uint64{1, 2, 3},
			index:       1,
			expected:    2,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := Value(tt.array, tt.index)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCheckSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		array       []uint64
		size        int
		expectError bool
	}{
		{
			name:        "nil array",
			array:       nil,
			size:        1,
			expectError: true,
		},
		{
			name:        "empty array",
			array:       []uint64{},
			size:        1,
			expectError: true,
		},
		{
			name:        "size too large",
			array:       []uint64{1, 2, 3},
			size:        4,
			expectError: true,
		},
		{
			name:        "valid size",
			array:       []uint64{1, 2, 3},
			size:        2,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := CheckSize(tt.array, tt.size)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDelElement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []uint64
		delId    uint64
		expected []uint64
	}{
		{
			name:     "empty array",
			array:    []uint64{},
			delId:    1,
			expected: []uint64{},
		},
		{
			name:     "element exists",
			array:    []uint64{1, 2, 3, 2, 4},
			delId:    2,
			expected: []uint64{1, 3, 4},
		},
		{
			name:     "element does not exist",
			array:    []uint64{1, 3, 4},
			delId:    2,
			expected: []uint64{1, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := DelElement(tt.array, tt.delId)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetNotZeroCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		array    []uint64
		expected uint64
	}{
		{
			name:     "empty array",
			array:    []uint64{},
			expected: 0,
		},
		{
			name:     "all zeros",
			array:    []uint64{0, 0, 0},
			expected: 0,
		},
		{
			name:     "mixed values",
			array:    []uint64{1, 0, 2, 0, 3},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := GetNotZeroCount(tt.array)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkContains(b *testing.B) {
	array := []uint64{1, 2, 3, 4, 5}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Contains(array, 3)
		}
	})
}

func BenchmarkIsElementRepeat(b *testing.B) {
	array := []uint64{1, 2, 2, 3, 4, 5}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IsElementRepeat(array, 0)
		}
	})
}

func BenchmarkIndex(b *testing.B) {
	array := []uint64{1, 2, 3, 4, 5}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Index(array, 3)
		}
	})
}

func BenchmarkCopy(b *testing.B) {
	array := []uint64{1, 2, 3, 4, 5}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Copy(array)
		}
	})
}

func BenchmarkRand(b *testing.B) {
	array := []uint64{1, 2, 3, 4, 5}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Rand(array, 3)
		}
	})
}

func BenchmarkDelElement(b *testing.B) {
	array := []uint64{1, 2, 3, 4, 5}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			b.StopTimer()

			testArray := Copy(array)

			b.StartTimer()
			DelElement(testArray, 3)
		}
	})
}
