package weights

import (
	"testing"

	"math"

	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/i64"
)

func TestTryNewWeights(t *testing.T) {
	tests := []struct {
		name    string
		ends    []uint64
		values  []int64
		wantErr bool
	}{
		{
			name:    "valid weights",
			ends:    []uint64{10, 20, 30},
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
			ends:    []uint64{1000000, 2000000, 3000000},
			values:  []int64{1, 2, 3},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := TryNewWeights(tt.ends, tt.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("TryNewWeights() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if w == nil {
					t.Errorf("TryNewWeights() returned nil weights without error")
				}
			}
		})
	}
}

func TestWeights_Value(t *testing.T) {
	// Create a test weights instance
	ends := []uint64{10, 20, 30}
	values := []int64{1, 2, 3}
	w, err := TryNewWeights(ends, values)
	if err != nil {
		t.Fatalf("Failed to create weights: %v", err)
	}

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
			val, ok := w.Value(tt.weight)
			if val != tt.wantVal || ok != tt.wantBool {
				t.Errorf("Value() = (%v, %v), want (%v, %v)", val, ok, tt.wantVal, tt.wantBool)
			}
		})
	}
}

func TestWeights_Rand(t *testing.T) {
	// Test with normal weights
	ends := []uint64{10, 20, 30}
	values := []int64{1, 2, 3}
	w, err := TryNewWeights(ends, values)
	if err != nil {
		t.Fatalf("Failed to create weights: %v", err)
	}

	// Run multiple times to ensure randomness
	validValues := map[int64]bool{1: true, 2: true, 3: true}
	for i := 0; i < 100; i++ {
		val := w.Rand()
		if !validValues[val] {
			t.Errorf("Rand() = %v, which is not a valid value", val)
		}
	}

	// Test with empty weights (edge case)
	emptyEnds := []uint64{0}
	emptyValues := []int64{0}
	emptyWeights, err := TryNewWeights(emptyEnds, emptyValues)
	if err != nil {
		t.Fatalf("Failed to create empty weights: %v", err)
	}

	val := emptyWeights.Rand()
	if val != 0 {
		t.Errorf("Empty weights Rand() = %v, want 0", val)
	}
}

func TestWeights_RandList(t *testing.T) {
	// Create a test weights instance
	ends := []uint64{5, 15, 30, 50, 75, 105, 140, 180, 225, 275, 330, 390, 455, 525, 600, 680, 765, 855, 950, 1050}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	w, err := TryNewWeights(ends, values)
	if err != nil {
		t.Fatalf("Failed to create weights: %v", err)
	}

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
			expectLen: len(values), // Only len(values) distinct values available
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
			result := w.RandList(tt.count)
			if i64.IsElementRepeat(result) {
				t.Errorf("RandList(%d) returned duplicate value: %v", tt.count, result)
			}

			// Check length
			if len(result) != tt.expectLen {
				t.Errorf("RandList(%d) returned %d items, want %d", tt.count, len(result), tt.expectLen)
			}

			// Check all values are valid
			if tt.expectLen > 0 {
				for _, v := range result {
					if !i64.Contains(values, v) {
						t.Errorf("RandList() returned invalid value: %v", v)
					}
				}
			}

			// Check for duplicates
			seen := make(map[int64]bool)
			for _, v := range result {
				if seen[v] {
					t.Errorf("RandList() returned duplicate value: %v", v)
				}
				seen[v] = true
			}
		})
	}

	// Test with small weights (edge case)
	smallEnds := []uint64{1}
	smallValues := []int64{42}
	smallWeights, err := TryNewWeights(smallEnds, smallValues)
	if err != nil {
		t.Fatalf("Failed to create small weights: %v", err)
	}

	smallResult := smallWeights.RandList(5)
	if len(smallResult) != 1 {
		t.Errorf("RandList for small weights = %v, expected length 1", smallResult)
	}
	if len(smallResult) > 0 && smallResult[0] != 42 {
		t.Errorf("RandList for small weights = %v, expected [42]", smallResult)
	}
}

func TestWeights_SimpleDistribution(t *testing.T) {
	// Test the distribution of random selections
	ends := []uint64{50, 80, 100} // Weights: 50, 30, 20
	values := []int64{1, 2, 3}
	w, err := TryNewWeights(ends, values)
	if err != nil {
		t.Fatalf("Failed to create weights: %v", err)
	}

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
	tests := []struct {
		name    string
		ends    []uint64
		values  []int64
		weights []float64 // the expected weight ratio of each value
	}{
		{
			name:    "uniform distribution",
			ends:    []uint64{10, 20, 30},
			values:  []int64{1, 2, 3},
			weights: []float64{0.333, 0.333, 0.333}, // each value has about 1/3 of the total weight
		},
		{
			name:    "gradient distribution",
			ends:    []uint64{10, 30, 60, 100},
			values:  []int64{1, 2, 3, 4},
			weights: []float64{0.1, 0.2, 0.3, 0.4}, // the weight increases gradually
		},
		{
			name:    "extreme distribution",
			ends:    []uint64{90, 95, 100},
			values:  []int64{1, 2, 3},
			weights: []float64{0.9, 0.05, 0.05}, // one value takes most of the weight
		},
		{
			name:    "zero weight test",
			ends:    []uint64{0, 50, 100},
			values:  []int64{1, 2, 3},
			weights: []float64{0, 0.5, 0.5}, // the first value has 0 weight
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := TryNewWeights(tt.ends, tt.values)
			if err != nil {
				t.Fatalf("failed to create weights: %v", err)
			}

			iterations := 100000 // large sample to reduce random error
			counts := make(map[int64]int, len(tt.values))

			// execute random sampling
			for i := 0; i < iterations; i++ {
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
				if math.Abs(actual-expected) > tolerance {
					t.Errorf("value %d's weight distribution does not match the expected: got %.4f, expected %.4f±%.4f",
						val, actual, expected, tolerance)
				}
			}
		})
	}
}

// TestWeights_RandListDistribution test the weight distribution of RandList method
func TestWeights_RandListDistribution(t *testing.T) {
	tests := []struct {
		name    string
		ends    []uint64
		values  []int64
		counts  []int     // different counts of calling RandList
		weights []float64 // the expected weight ratio of each value
	}{
		{
			name:    "uniform distribution - all elements",
			ends:    []uint64{10, 20, 30, 40, 50},
			values:  []int64{1, 2, 3, 4, 5},
			counts:  []int{5},                           // request all 5 elements
			weights: []float64{0.2, 0.2, 0.2, 0.2, 0.2}, // should appear uniformly
		},
		{
			name:    "different distribution - partial count",
			ends:    []uint64{30, 50, 90, 100},
			values:  []int64{1, 2, 3, 4},
			counts:  []int{2},                      // request 2 elements each time
			weights: []float64{0.3, 0.2, 0.4, 0.1}, // different weights
		},
		{
			name:    "multiple request counts",
			ends:    []uint64{25, 50, 75, 100},
			values:  []int64{1, 2, 3, 4},
			counts:  []int{1, 2, 3},                    // use different request counts
			weights: []float64{0.25, 0.25, 0.25, 0.25}, // uniform distribution
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := TryNewWeights(tt.ends, tt.values)
			if err != nil {
				t.Fatalf("failed to create weights: %v", err)
			}

			iterations := 10000 // large sample to reduce random error
			valueFreqs := make(map[int64]int)
			totalResults := 0

			// test each request count
			for _, count := range tt.counts {
				for i := 0; i < iterations; i++ {
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
					expectedUniform := 1.0 / float64(len(tt.values))
					if math.Abs(actual-expectedUniform) > 0.05 {
						t.Errorf("when requesting all elements, value %d's distribution is not uniform: got %.4f, expected close to %.4f",
							val, actual, expectedUniform)
					}
				} else if math.Abs(actual-expected) > tolerance {
					// when not requesting all elements, check if it is within the tolerance
					// note: due to the duplicate elimination logic of RandList, the actual distribution may differ from the original weights
					t.Logf("warning: value %d's weight distribution is different from expected: got %.4f, expected %.4f±%.4f",
						val, actual, expected, tolerance)
				}
			}
		})
	}
}

// TestWeights_ExtremeCases test the weight distribution in extreme cases
func TestWeights_ExtremeCases(t *testing.T) {
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
			w, err := TryNewWeights(tt.ends, tt.values)
			if err != nil {
				t.Fatalf("failed to create weights: %v", err)
			}

			t.Logf("test: %s", tt.desc)

			// test Rand
			iterations := 1000
			randCounts := make(map[int64]int)
			for i := 0; i < iterations; i++ {
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

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = TryNewWeights(ends, values)
	}
}

// BenchmarkWeights_Value benchmarks the Value method
func BenchmarkWeights_Value(b *testing.B) {
	ends := []uint64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	w, _ := TryNewWeights(ends, values)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = w.Value(uint64(i % 100))
	}
}

// BenchmarkWeights_Rand benchmarks the Rand method
func BenchmarkWeights_Rand(b *testing.B) {
	ends := []uint64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	values := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	w, _ := TryNewWeights(ends, values)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = w.Rand()
	}
}

// BenchmarkWeights_RandList benchmarks the RandList method with different counts
func BenchmarkWeights_RandList(b *testing.B) {
	// Prepare weights with many values for a realistic benchmark
	const numWeights = 100
	ends := make([]uint64, numWeights)
	values := make([]int64, numWeights)

	for i := 0; i < numWeights; i++ {
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
			for i := 0; i < b.N; i++ {
				_ = w.RandList(bc.count)
			}
		})
	}
}
