package weights

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTryNewWeights(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ends    []uint64
		values  []int64
		wantErr bool
	}{
		{
			name:    "valid weights",
			ends:    []uint64{10, 10, 10},
			values:  []int64{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "empty ends",
			ends:    []uint64{},
			values:  []int64{},
			wantErr: true,
		},
		{
			name:    "length mismatch",
			ends:    []uint64{10, 20},
			values:  []int64{1},
			wantErr: true,
		},
		{
			name:    "zero weight",
			ends:    []uint64{0, 10},
			values:  []int64{1, 2},
			wantErr: false,
		},
		{
			name:    "large weights",
			ends:    []uint64{1000000, 1000000, 1000000},
			values:  []int64{1, 2, 3},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w, err := TryNewWeights(tt.ends, tt.values)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.wantErr, w == nil)
		})
	}
}

func TestWeights_Value(t *testing.T) {
	t.Parallel()

	// Create a test weights instance
	ends := []uint64{10, 10, 10}
	values := []int64{1, 2, 3}

	w, err := TryNewWeights(ends, values)
	require.NoError(t, err)

	tests := []struct {
		name     string
		weight   uint64
		wantVal  int64
		wantBool bool
	}{
		{
			name:     "first range",
			weight:   5,
			wantVal:  1,
			wantBool: true,
		},
		{
			name:     "second range",
			weight:   15,
			wantVal:  2,
			wantBool: true,
		},
		{
			name:     "third range",
			weight:   25,
			wantVal:  3,
			wantBool: true,
		},
		{
			name:     "boundary - first range",
			weight:   0,
			wantVal:  1,
			wantBool: true,
		},
		{
			name:     "boundary - between first and second",
			weight:   10,
			wantVal:  2,
			wantBool: true,
		},
		{
			name:     "boundary - second and third",
			weight:   20,
			wantVal:  3,
			wantBool: true,
		},
		{
			name:     "out of range - too large",
			weight:   35,
			wantVal:  0,
			wantBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			val, ok := w.Value(tt.weight)
			assert.Equal(t, tt.wantVal, val)
			assert.Equal(t, tt.wantBool, ok)
		})
	}
}

func TestWeights_Rand(t *testing.T) {
	t.Parallel()

	// Test with normal weights
	ends := []uint64{10, 10, 10}
	values := []int64{1, 2, 3}

	w, err := TryNewWeights(ends, values)
	require.NoError(t, err)

	// Run multiple times to ensure randomness
	validValues := map[int64]bool{1: true, 2: true, 3: true}

	for range 100 {
		val := w.Rand()
		if !validValues[val] {
			t.Errorf("Rand() = %v, which is not a valid value", val)
		}
	}

	// Test with empty weights (edge case)
	emptyEnds := []uint64{0}
	emptyValues := []int64{0}

	emptyWeights, err := TryNewWeights(emptyEnds, emptyValues)
	require.NoError(t, err)

	val := emptyWeights.Rand()
	assert.Equal(t, int64(0), val)
}

func TestWeights_RandList(t *testing.T) {
	t.Parallel()

	// Create a test weights instance
	ends := []uint64{5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 65, 70, 75, 80, 85, 90, 95, 100}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

	w, err := TryNewWeights(ends, values)
	require.NoError(t, err)

	tests := []struct {
		name  string
		count int
		// We can't predict exact values due to randomness,
		// but we can check length and value range
		expectLen int
	}{
		{
			name:      "normal case",
			count:     3,
			expectLen: 3,
		},
		{
			name:      "max count case",
			count:     len(values),
			expectLen: len(values),
		},
		{
			name:      "request more than available",
			count:     len(values) + 5,
			expectLen: len(values) + 5,
		},
		{
			name:      "zero count",
			count:     0,
			expectLen: 0,
		},
		{
			name:      "negative count",
			count:     -1,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := w.RandList(tt.count)

			assert.Equal(t, tt.expectLen, len(result))

			if tt.expectLen > 0 {
				for _, v := range result {
					assert.Contains(t, values, v)
				}
			}
		})
	}
}

func TestWeights_RandList_EdgeCases(t *testing.T) {
	t.Parallel()

	// Test with small weights (edge case)
	ends := []uint64{1}
	values := []int64{42}

	weights, err := TryNewWeights(ends, values)
	require.NoError(t, err)

	ret := weights.RandList(5)
	assert.Equal(t, 5, len(ret))

	for _, v := range ret {
		assert.Contains(t, values, v)
	}
}

func TestWeights_SimpleDistribution(t *testing.T) {
	t.Parallel()

	// Test the distribution of random selections
	ends := []uint64{50, 30, 20} // Weights: 50, 30, 20
	values := []int64{1, 2, 3}

	w, err := TryNewWeights(ends, values)
	require.NoError(t, err)

	counts := make(map[int64]int)
	iterations := 10000

	for i := 0; i < iterations; i++ {
		val := w.Rand()
		counts[val]++
	}

	// Check if distribution roughly matches weights
	// 50% for value 1, 30% for value 2, 20% for value 3
	expectedRatios := map[int64]float64{
		1: 0.5,
		2: 0.3,
		3: 0.2,
	}

	t.Logf("Distribution after %d iterations:", iterations)

	for val, count := range counts {
		ratio := float64(count) / float64(iterations)
		expected := expectedRatios[val]
		t.Logf("Value %d: count=%d, ratio=%.4f, expected=%.4f", val, count, ratio, expected)

		// Allow for some statistical variation (±5%)
		if ratio < expected-0.05 || ratio > expected+0.05 {
			t.Errorf("Distribution for value %d: got ratio %.4f, expected around %.4f",
				val, ratio, expected)
		}
	}
}

// TestWeights_DetailedRandDistribution test the weight distribution of Rand method
func TestWeights_DetailedRandDistribution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ends    []uint64
		values  []int64
		weights []float64 // the expected weight ratio of each value
	}{
		{
			name:    "uniform distribution",
			ends:    []uint64{10, 10, 10},
			values:  []int64{1, 2, 3},
			weights: []float64{0.333, 0.333, 0.333}, // each value has about 1/3 of the total weight
		},
		{
			name:    "gradient distribution",
			ends:    []uint64{10, 20, 30, 40},
			values:  []int64{1, 2, 3, 4},
			weights: []float64{0.1, 0.2, 0.3, 0.4}, // the weight increases gradually
		},
		{
			name:    "extreme distribution",
			ends:    []uint64{90, 5, 5},
			values:  []int64{1, 2, 3},
			weights: []float64{0.9, 0.05, 0.05}, // one value takes most of the weight
		},
		{
			name:    "zero weight test",
			ends:    []uint64{0, 50, 50},
			values:  []int64{1, 2, 3},
			weights: []float64{0, 0.5, 0.5}, // the first value has 0 weight
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w, err := TryNewWeights(tt.ends, tt.values)
			if err != nil {
				t.Fatalf("failed to create weights: %v", err)
			}

			iterations := 100000 // large sample to reduce random error
			counts := make(map[int64]int, len(tt.values))

			// execute random sampling
			for range iterations {
				val := w.Rand()
				counts[val]++
			}

			// verify the distribution
			t.Logf("weight distribution test '%s' (%d iterations):", tt.name, iterations)

			for i, val := range tt.values {
				expected := tt.weights[i]
				count := counts[val]
				actual := float64(count) / float64(iterations)
				tolerance := 0.02 // allow 2% error range

				t.Logf("value %d: count=%d, ratio=%.4f, expected=%.4f, error=%.4f",
					val, count, actual, expected, math.Abs(actual-expected))

				// check if it is within the tolerance
				assert.True(t, math.Abs(actual-expected) <= tolerance)
			}
		})
	}
}

// TestWeights_RandListDistribution test the weight distribution of RandList method
func TestWeights_RandListDistribution(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ends    []uint64
		values  []int64
		counts  []int     // different counts of calling RandList
		weights []float64 // the expected weight ratio of each value
	}{
		{
			name:    "uniform distribution - all elements",
			ends:    []uint64{10, 10, 10, 10, 10},
			values:  []int64{1, 2, 3, 4, 5},
			counts:  []int{5},                           // request all 5 elements
			weights: []float64{0.2, 0.2, 0.2, 0.2, 0.2}, // should appear uniformly
		},
		{
			name:    "different distribution - partial count",
			ends:    []uint64{30, 20, 40, 10},
			values:  []int64{1, 2, 3, 4},
			counts:  []int{2},                      // request 2 elements each time
			weights: []float64{0.3, 0.2, 0.4, 0.1}, // different weights
		},
		{
			name:    "multiple request counts",
			ends:    []uint64{25, 25, 25, 25},
			values:  []int64{1, 2, 3, 4},
			counts:  []int{1, 2, 3},                    // use different request counts
			weights: []float64{0.25, 0.25, 0.25, 0.25}, // uniform distribution
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w, err := TryNewWeights(tt.ends, tt.values)
			require.NoError(t, err)

			iterations := 10000 // large sample to reduce random error
			valueFreqs := make(map[int64]int)
			totalResults := 0

			// test each request count
			for _, count := range tt.counts {
				for range iterations {
					results := w.RandList(count)
					for _, val := range results {
						valueFreqs[val]++
						totalResults++
					}
				}
			}

			// analyze the results
			t.Logf("RandList distribution test '%s' (total results: %d):", tt.name, totalResults)

			// calculate the expected frequency of each value
			for i, val := range tt.values {
				expected := tt.weights[i]
				freq := valueFreqs[val]
				actual := float64(freq) / float64(totalResults)

				// RandList's characteristics determine that it may differ significantly from pure randomness
				// so use a more lenient tolerance
				tolerance := 0.05

				t.Logf("value %d: count=%d, ratio=%.4f, expected=%.4f, error=%.4f",
					val, freq, actual, expected, math.Abs(actual-expected))

				// for the case of full request, the frequency should be closer to uniform distribution
				if len(tt.counts) == 1 && tt.counts[0] == len(tt.values) {
					assert.True(t, math.Abs(actual-(1.0/float64(len(tt.values)))) <= 0.05)
				} else {
					assert.True(t, math.Abs(actual-expected) <= tolerance)
				}
			}
		})
	}
}

// TestWeights_ExtremeCases test the weight distribution in extreme cases
func TestWeights_ExtremeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		ends   []uint64
		values []int64
		desc   string // test description
	}{
		{
			name:   "single element",
			ends:   []uint64{100},
			values: []int64{42},
			desc:   "test the random selection when there is only one element",
		},
		{
			name:   "maximum weight difference",
			ends:   []uint64{1, 10000000},
			values: []int64{1, 2},
			desc:   "test the randomness when the weights are extremely unbalanced",
		},
		{
			name:   "zero weight element",
			ends:   []uint64{0, 0, 100},
			values: []int64{1, 2, 3},
			desc:   "test the handling of elements with zero weight",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			w, err := TryNewWeights(tt.ends, tt.values)
			require.NoError(t, err)

			t.Logf("test: %s", tt.desc)

			// test Rand
			iterations := 1000
			randCounts := make(map[int64]int)

			for range iterations {
				val := w.Rand()
				randCounts[val]++
			}

			// output Rand results
			t.Logf("Rand results distribution (%d iterations):", iterations)

			for val, count := range randCounts {
				ratio := float64(count) / float64(iterations)
				t.Logf("value %d: count=%d, ratio=%.4f", val, count, ratio)
			}

			// test RandList
			if len(tt.values) > 1 {
				randListResults := w.RandList(len(tt.values) - 1)
				t.Logf("RandList(%d) results: %v", len(tt.values)-1, randListResults)
			}
		})
	}
}

// BenchmarkTryNewWeights benchmarks the creation of a new Weights instance
func BenchmarkTryNewWeights(b *testing.B) {
	ends := []uint64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = TryNewWeights(ends, values)
		}
	})
}

// BenchmarkWeights_Value benchmarks the Value method
func BenchmarkWeights_Value(b *testing.B) {
	ends := []uint64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	w, _ := TryNewWeights(ends, values)

	b.RunParallel(func(pb *testing.PB) {
		var i uint64
		for pb.Next() {
			_, _ = w.Value(i % 100)
			i++
		}
	})
}

// BenchmarkWeights_Rand benchmarks the Rand method
func BenchmarkWeights_Rand(b *testing.B) {
	ends := []uint64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	w, _ := TryNewWeights(ends, values)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = w.Rand()
		}
	})
}

// BenchmarkWeights_RandList benchmarks the RandList method with different counts
func BenchmarkWeights_RandList(b *testing.B) {
	// Prepare weights with many values for a realistic benchmark
	const numWeights = 100
	ends := make([]uint64, numWeights)
	values := make([]int64, numWeights)

	for i := range numWeights {
		ends[i] = uint64((i + 1) * 10)
		values[i] = int64(i + 1)
	}

	w, _ := TryNewWeights(ends, values)

	benchCases := []struct {
		name  string
		count int
	}{
		{"Small", 5},
		{"Medium", 25},
		{"Large", 50},
		{"ExtraLarge", 90},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = w.RandList(bc.count)
				}
			})
		})
	}
}
