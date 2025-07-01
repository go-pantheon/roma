package u64

import (
	"math"
	"math/rand/v2"
	"strconv"
)

func Max(x, y uint64) uint64 {
	if x > y {
		return x
	}

	return y
}

func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}

	return y
}

func Sub(x, y uint64) uint64 {
	if x <= y {
		return 0
	}

	return x - y
}

// Add performs saturating addition of two uint64 numbers.
// If adding two positive numbers results in a negative (overflow), returns MaxUint64.
// If adding two negative numbers results in a positive (underflow), returns MinUint64.
// Otherwise returns the normal sum.
func Add(x, y uint64) uint64 {
	r := x + y
	if r < x || r < y {
		return math.MaxUint64
	}

	return r
}

func Random(v uint64) uint64 {
	if v <= 0 {
		return 0
	}

	return rand.Uint64N(v)
}

func Divide2f64(x, y uint64) float64 {
	if x == 0 || y == 0 {
		return 0
	}

	return float64(x) / float64(y)
}

func F64WithDigits(v float64, dig int) float64 {
	scale := math.Pow10(dig)
	return math.Round(v*scale) / scale
}

func Pow(a, n uint64) uint64 {
	var x uint64 = 1

	for n != 0 {
		if n&1 == 1 {
			x *= a
		}

		n >>= 1
		a *= a
	}

	return x
}

func Exp(x uint64) uint64 {
	return uint64(math.Exp(float64(x)))
}

func CeilDivide(x, y uint64) uint64 {
	if x == 0 || y == 0 {
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

func ToU32(x uint64) uint32 {
	if x <= 0 {
		return 0
	}

	if x > math.MaxUint32 {
		return math.MaxUint32
	}

	return uint32(x)
}

func ToU64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
