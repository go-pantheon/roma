package i64

import (
	"math"
	"math/rand"
	"strconv"
)

func Max(x, y int64) int64 {
	if x > y {
		return x
	}

	return y
}

func Min(x, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

// Reduce returns the difference between x and y when x > y, otherwise returns 0
// Example: Reduce(5,3) = 2; Reduce(3,5) = 0
func Reduce(x, y int64) int64 {
	if x <= y {
		return 0
	}

	return x - y
}

// Add performs safe addition of two int64 numbers, handling overflow cases.
// If adding two positive numbers results in a negative (overflow), returns MaxInt64.
// If adding two negative numbers results in a positive (underflow), returns MinInt64.
// Otherwise returns the normal sum.
func Add(x, y int64) int64 {
	if x == 0 {
		return y
	}
	
	if y == 0 {
		return x
	}

	r := x + y

	// Check for positive overflow
	if x > 0 && y > 0 && r < 0 {
		return math.MaxInt64
	}

	// Check for negative overflow
	if x < 0 && y < 0 && r > 0 {
		return math.MinInt64
	}

	return r
}

// Random generates a random number between 0 and v-1
// Note: Call rand.Seed() before using this function
// to ensure proper randomness
func Random(v int64) int64 {
	if v == 0 {
		return 0
	}

	return rand.Int63n(v)
}

func Divide2f64(x, y int64) float64 {
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
func Pow(a, n int64) int64 {
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
	var result int64 = 1

	for n > 0 {
		if n&1 == 1 {
			// Check for overflow before multiplying
			if willOverflow(result, a) {
				return math.MaxInt64
			}

			result *= a
		}

		n >>= 1
		if n != 0 { // Skip last squaring
			if a > math.MaxInt64/a { // check overflow
				return math.MaxInt64
			}

			a *= a
		}
	}

	return result
}

// willOverflow checks if multiplying a and b will overflow int64
func willOverflow(a, b int64) bool {
	if a == 0 || b == 0 {
		return false
	}

	result := a * b

	return a != result/b
}

func Exp(x int64) int64 {
	return int64(math.Exp(float64(x)))
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}

	return x
}

func CeilDivide(x, y int64) int64 {
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

func ToI32(x int64) int32 {
	if x <= 0 {
		return 0
	}

	if x > math.MaxInt32 {
		return math.MaxInt32
	}
	
	return int32(x)
}

func ToI64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
