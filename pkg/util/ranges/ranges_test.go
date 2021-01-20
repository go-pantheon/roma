package ranges

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryNewPair(t *testing.T) {
	tests := []struct {
		name      string
		start     uint64
		end       uint64
		value     int64
		expectErr bool
	}{
		{
			name:      "valid range",
			start:     10,
			end:       20,
			value:     5,
			expectErr: false,
		},
		{
			name:      "valid range with same start and end",
			start:     10,
			end:       10,
			value:     5,
			expectErr: false,
		},
		{
			name:      "invalid range with start > end",
			start:     20,
			end:       10,
			value:     5,
			expectErr: true,
		},
		{
			name:      "valid range with zero value",
			start:     10,
			end:       20,
			value:     0,
			expectErr: false,
		},
		{
			name:      "valid range with max uint64 value",
			start:     0,
			end:       ^uint64(0),
			value:     5,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pair, err := TryNewPair(tt.start, tt.end, tt.value)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, pair)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pair)
				assert.Equal(t, tt.start, pair.Start)
				assert.Equal(t, tt.end, pair.End)
				assert.Equal(t, tt.value, pair.Value)
			}
		})
	}
}

func TestPair_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		pair     *Pair
		expected bool
	}{
		{
			name:     "valid range",
			pair:     &Pair{Start: 10, End: 20, Value: 5},
			expected: true,
		},
		{
			name:     "valid range with same start and end",
			pair:     &Pair{Start: 10, End: 10, Value: 5},
			expected: true,
		},
		{
			name:     "invalid range with start > end",
			pair:     &Pair{Start: 20, End: 10, Value: 5},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.pair.IsValid())
		})
	}
}

func TestPair_Contains(t *testing.T) {
	pair := &Pair{Start: 10, End: 20, Value: 5}

	tests := []struct {
		name     string
		value    uint64
		expected bool
	}{
		{
			name:     "value at start boundary",
			value:    10,
			expected: true,
		},
		{
			name:     "value at end boundary",
			value:    20,
			expected: false,
		},
		{
			name:     "value within range",
			value:    15,
			expected: true,
		},
		{
			name:     "value before range",
			value:    5,
			expected: false,
		},
		{
			name:     "value after range",
			value:    25,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, pair.Contains(tt.value))
		})
	}
}

func TestTryNewRange(t *testing.T) {
	tests := []struct {
		name      string
		coords    []uint64
		values    []int64
		expectErr bool
	}{
		{
			name:      "valid ranges",
			coords:    []uint64{10, 20, 30, 40, 50, 60},
			values:    []int64{1, 2, 3},
			expectErr: false,
		},
		{
			name:      "empty coords",
			coords:    []uint64{},
			values:    []int64{},
			expectErr: true,
		},
		{
			name:      "odd number of coords",
			coords:    []uint64{10, 20, 30},
			values:    []int64{1, 2},
			expectErr: true,
		},
		{
			name:      "mismatched lengths",
			coords:    []uint64{10, 20, 30, 40},
			values:    []int64{1},
			expectErr: true,
		},
		{
			name:      "invalid range pair",
			coords:    []uint64{20, 10, 30, 40},
			values:    []int64{1, 2},
			expectErr: true,
		},
		{
			name:      "overlapping ranges",
			coords:    []uint64{10, 30, 20, 40},
			values:    []int64{1, 2},
			expectErr: true,
		},
		{
			name:      "valid unsorted ranges (should sort internally)",
			coords:    []uint64{30, 40, 10, 20, 50, 60},
			values:    []int64{2, 1, 3},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := TryNewRange(tt.coords, tt.values)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, len(tt.values), len(r.Pairs))

				// Check if pairs are sorted by start
				for i := 1; i < len(r.Pairs); i++ {
					assert.Less(t, r.Pairs[i-1].Start, r.Pairs[i].Start)
				}

				// Check if pairs are not overlapping
				for i := 1; i < len(r.Pairs); i++ {
					assert.LessOrEqual(t, r.Pairs[i-1].End, r.Pairs[i].Start)
				}
			}
		})
	}
}

func TestRange_Find(t *testing.T) {
	// Create a test range
	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	r, err := TryNewRange(coords, values)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	tests := []struct {
		name        string
		value       uint64
		expectedVal int64
		expectedIdx int
	}{
		{
			name:        "value in first range",
			value:       15,
			expectedVal: 1,
			expectedIdx: 0,
		},
		{
			name:        "value at start of first range",
			value:       10,
			expectedVal: 1,
			expectedIdx: 0,
		},
		{
			name:        "value at end of first range",
			value:       20,
			expectedVal: 0,
			expectedIdx: -1,
		},
		{
			name:        "value in second range",
			value:       35,
			expectedVal: 2,
			expectedIdx: 1,
		},
		{
			name:        "value in last range",
			value:       55,
			expectedVal: 3,
			expectedIdx: 2,
		},
		{
			name:        "value before all ranges",
			value:       5,
			expectedVal: 0,
			expectedIdx: -1,
		},
		{
			name:        "value after all ranges",
			value:       65,
			expectedVal: 0,
			expectedIdx: -1,
		},
		{
			name:        "value in gap between ranges",
			value:       25,
			expectedVal: 0,
			expectedIdx: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, idx := r.Find(tt.value)
			assert.Equal(t, tt.expectedVal, val)
			assert.Equal(t, tt.expectedIdx, idx)
		})
	}

	// Test with empty range
	emptyRange := &Range{Pairs: []*Pair{}}
	val, idx := emptyRange.Find(10)
	assert.Equal(t, int64(0), val)
	assert.Equal(t, -1, idx)
}

func TestRange_MinMax(t *testing.T) {
	// Test with regular range
	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	r, err := TryNewRange(coords, values)
	assert.NoError(t, err)

	assert.Equal(t, uint64(10), r.Min())
	assert.Equal(t, uint64(60), r.Max())

	// Test with single range
	singleCoords := []uint64{10, 20}
	singleValues := []int64{1}
	singleRange, err := TryNewRange(singleCoords, singleValues)
	assert.NoError(t, err)

	assert.Equal(t, uint64(10), singleRange.Min())
	assert.Equal(t, uint64(20), singleRange.Max())

	// Test with empty range
	emptyRange := &Range{Pairs: []*Pair{}}
	assert.Equal(t, uint64(0), emptyRange.Min())
	assert.Equal(t, uint64(0), emptyRange.Max())
}

func TestRange_Len(t *testing.T) {
	// Test with regular range
	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	r, err := TryNewRange(coords, values)
	assert.NoError(t, err)

	assert.Equal(t, 3, r.Len())

	// Test with single range
	singleCoords := []uint64{10, 20}
	singleValues := []int64{1}
	singleRange, err := TryNewRange(singleCoords, singleValues)
	assert.NoError(t, err)

	assert.Equal(t, 1, singleRange.Len())

	// Test with empty range
	emptyRange := &Range{Pairs: []*Pair{}}
	assert.Equal(t, 0, emptyRange.Len())
}

func TestRange_Rand(t *testing.T) {
	// Mock random to make test deterministic
	rand.Seed(42)

	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	r, err := TryNewRange(coords, values)
	assert.NoError(t, err)

	// Test that Rand returns a value from one of the pairs
	for i := 0; i < 100; i++ {
		val, idx := r.Rand()
		assert.GreaterOrEqual(t, idx, 0)
		assert.Less(t, idx, len(r.Pairs))
		assert.Equal(t, r.Pairs[idx].Value, val)
	}

	// Test with empty range
	emptyRange := &Range{Pairs: []*Pair{}}
	val, idx := emptyRange.Rand()
	assert.Equal(t, int64(0), val)
	assert.Equal(t, 0, idx)
}

func TestRange_validate(t *testing.T) {
	tests := []struct {
		name      string
		pairs     []*Pair
		expectErr bool
	}{
		{
			name: "valid ranges",
			pairs: []*Pair{
				{Start: 10, End: 20, Value: 1},
				{Start: 30, End: 40, Value: 2},
				{Start: 50, End: 60, Value: 3},
			},
			expectErr: false,
		},
		{
			name: "invalid range (start > end)",
			pairs: []*Pair{
				{Start: 10, End: 20, Value: 1},
				{Start: 50, End: 40, Value: 2},
			},
			expectErr: true,
		},
		{
			name: "overlapping ranges",
			pairs: []*Pair{
				{Start: 10, End: 30, Value: 1},
				{Start: 20, End: 40, Value: 2},
			},
			expectErr: true,
		},
		{
			name: "adjacent ranges (not overlapping)",
			pairs: []*Pair{
				{Start: 10, End: 20, Value: 1},
				{Start: 20, End: 30, Value: 2},
			},
			expectErr: false,
		},
		{
			name:      "empty ranges",
			pairs:     []*Pair{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Range{Pairs: tt.pairs}
			err := r.validate()
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Benchmarks

func BenchmarkTryNewPair(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = TryNewPair(10, 20, 5)
	}
}

func BenchmarkPair_Contains(b *testing.B) {
	pair := &Pair{Start: 10, End: 20, Value: 5}
	for i := 0; i < b.N; i++ {
		_ = pair.Contains(15)
	}
}

func BenchmarkTryNewRange(b *testing.B) {
	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TryNewRange(coords, values)
	}
}

func BenchmarkRange_Find(b *testing.B) {
	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	r, _ := TryNewRange(coords, values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Find(35)
	}
}

func BenchmarkRange_Find_Large(b *testing.B) {
	// Create a large range with 1000 pairs
	coords := make([]uint64, 2000)
	values := make([]int64, 1000)

	for i := 0; i < 1000; i++ {
		coords[i*2] = uint64(i * 10)
		coords[i*2+1] = uint64(i*10 + 5)
		values[i] = int64(i)
	}

	r, _ := TryNewRange(coords, values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Find(uint64(rand.Intn(10000)))
	}
}

func BenchmarkRange_Rand(b *testing.B) {
	coords := []uint64{10, 20, 30, 40, 50, 60}
	values := []int64{1, 2, 3}
	r, _ := TryNewRange(coords, values)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = r.Rand()
	}
}
