// Package f64 provides basic mathematical operations for float64 numbers.
package f64

import "math"

// Max returns the larger of x or y.
// Special cases are:
//
//	Max(x, NaN) = NaN
//	Max(NaN, x) = NaN
func Max(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}

	if x > y {
		return x
	}

	return y
}

// Min returns the smaller of x or y.
// Special cases are:
//
//	Min(x, NaN) = NaN
//	Min(NaN, x) = NaN
func Min(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}

	if x < y {
		return x
	}

	return y
}

// Reduce subtracts y from x if x is greater than y, otherwise returns 0.
// Special cases are:
//
//	Reduce(x, NaN) = NaN
//	Reduce(NaN, x) = NaN
//	Reduce(±Inf, ±Inf) = NaN
func Reduce(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}

	if x > y {
		result := x - y
		// minimum denormalized float64
		if result != 0 && math.Abs(result) < 2.2250738585072014e-308 {
			result = math.Nextafter(result, math.Inf(int(math.Copysign(1, result))))
		}

		return result
	}

	return 0
}

// Add returns the sum of x and y, capped at math.MaxFloat64.
// Special cases are:
//
//	Add(x, NaN) = NaN
//	Add(NaN, x) = NaN
//	Add(+Inf, -Inf) = NaN
//	Add(-Inf, +Inf) = NaN
func Add(x, y float64) float64 {
	if math.IsNaN(x) || math.IsNaN(y) {
		return math.NaN()
	}

	// Handle opposite infinity cases first
	if (math.IsInf(x, 1) && math.IsInf(y, -1)) || (math.IsInf(x, -1) && math.IsInf(y, 1)) {
		return math.NaN()
	}

	// Handle same infinity cases
	if math.IsInf(x, 1) || math.IsInf(y, 1) {
		return math.Inf(1)
	}

	if math.IsInf(x, -1) || math.IsInf(y, -1) {
		return math.Inf(-1)
	}

	sum := x + y

	// Handle overflow to infinity cases
	if math.IsInf(sum, 1) {
		return math.MaxFloat64
	}

	if math.IsInf(sum, -1) {
		return -math.MaxFloat64
	}

	// Cap finite values at MaxFloat64 boundaries
	if sum > math.MaxFloat64 {
		return math.MaxFloat64
	}

	if sum < -math.MaxFloat64 {
		return -math.MaxFloat64
	}

	return sum
}
