package i32

import (
	"math"
	"math/rand/v2"
)

func Max(x, y int32) int32 {
	if x > y {
		return x
	}

	return y
}

func Min(x, y int32) int32 {
	if x < y {
		return x
	}

	return y
}

// Reduce returns the difference between x and y when x > y, otherwise returns 0
// Example: Reduce(5,3) = 2; Reduce(3,5) = 0
func Reduce(x, y int32) int32 {
	if x <= y {
		return 0
	}

	return x - y
}

// Add performs safe addition of two int32 numbers, handling overflow cases.
// If adding two positive numbers results in a negative (overflow), returns MaxInt32.
// If adding two negative numbers results in a positive (underflow), returns MinInt64.
// Otherwise returns the normal sum.
func Add(x, y int32) int32 {
	if x == 0 {
		return y
	}

	if y == 0 {
		return x
	}

	r := x + y

	// Check for positive overflow
	if x > 0 && y > 0 && r < 0 {
		return math.MaxInt32
	}

	// Check for negative overflow
	if x < 0 && y < 0 && r > 0 {
		return math.MinInt32
	}

	return r
}

// Random generates a random number between 0 and v-1
// Note: Call rand.Seed() before using this function
// to ensure proper randomness
func Random(v int32) int32 {
	if v <= 0 {
		return 0
	}

	return rand.Int32N(v)
}

func Divide2f64(x, y int32) float64 {
	if x == 0 || y == 0 {
		return 0
	}

	return float64(x) / float64(y)
}

func F64WithDigits(v float64, dig int) float64 {
	shift := math.Pow(10, float64(dig))
	return math.Round(v*shift) / shift
}

// Pow calculates a^n using binary exponentiation algorithm.
// Returns 0 for negative exponents and 1 for n=0.
// Note: Does not check for integer overflow.
func Pow(a, n int32) int32 {
	// Handle edge cases
	if n < 0 {
		return 0 // Return 0 for negative exponents
	}

	if n == 0 {
		return 1 // Any number to power 0 is 1
	}

	if a == 0 {
		return 0 // 0 to any positive power is 0
	}

	if a == 1 {
		return 1 // 1 to any power is 1
	}

	// Binary exponentiation
	var result int32 = 1

	for n > 0 {
		if n&1 == 1 {
			// Check for overflow before multiplying
			if willOverflow(result, a) {
				return math.MaxInt32
			}

			result *= a
		}

		n >>= 1
		if n != 0 { // Skip last squaring
			if a > math.MaxInt32/a { // check overflow
				return math.MaxInt32
			}

			a *= a
		}
	}

	return result
}

// willOverflow checks if multiplying a and b will overflow int32
func willOverflow(a, b int32) bool {
	if a == 0 || b == 0 {
		return false
	}

	result := a * b

	return a != result/b
}

func Exp(x int32) int32 {
	return int32(math.Exp(float64(x)))
}

func Abs(x int32) int32 {
	if x < 0 {
		return -x
	}

	return x
}

func CeilDivide(x, y int32) int32 {
	if x == 0 || y == 0 {
		return 0
	}

	if x >= 0 && y < 0 {
		return 0
	}

	if x <= y {
		return 1
	}

	if x%y > 0 {
		return x/y + 1
	}

	return x / y
}
