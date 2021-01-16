package u64

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		x, y     uint64
		expected uint64
	}{
		{"x greater than y", 10, 5, 10},
		{"y greater than x", 5, 10, 10},
		{"equal values", 5, 5, 5},
		{"zero values", 0, 0, 0},
		{"max uint64", math.MaxUint64, 5, math.MaxUint64},
		{"both max uint64", math.MaxUint64, math.MaxUint64, math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Max(tt.x, tt.y))
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		x, y     uint64
		expected uint64
	}{
		{"x less than y", 5, 10, 5},
		{"y less than x", 10, 5, 5},
		{"equal values", 5, 5, 5},
		{"zero values", 0, 0, 0},
		{"one zero", 0, 5, 0},
		{"max uint64", math.MaxUint64, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Min(tt.x, tt.y))
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		x, y     uint64
		expected uint64
	}{
		{"normal reduction", 10, 3, 7},
		{"equal values", 5, 5, 0},
		{"x less than y", 3, 5, 0},
		{"zero values", 0, 0, 0},
		{"reduce from max", math.MaxUint64, 1, math.MaxUint64 - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Sub(tt.x, tt.y))
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		x, y     uint64
		expected uint64
	}{
		{"normal addition", 5, 3, 8},
		{"zero addition", 5, 0, 5},
		{"overflow case", math.MaxUint64, 1, math.MaxUint64},
		{"both max", math.MaxUint64, math.MaxUint64, math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Add(tt.x, tt.y))
		})
	}
}

func TestRandom(t *testing.T) {
	tests := []struct {
		name string
		v    uint64
	}{
		{"zero value", 0},
		{"small value", 10},
		{"large value", math.MaxInt64},
		{"overflow value", math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Random(tt.v)
			if tt.v == 0 {
				assert.Equal(t, uint64(0), got)
			} else {
				assert.True(t, got < tt.v)
			}
		})
	}
}

func TestDivide2f64(t *testing.T) {
	tests := []struct {
		name     string
		x, y     uint64
		expected float64
	}{
		{"normal division", 10, 2, 5.0},
		{"decimal result", 5, 2, 2.5},
		{"x is zero", 0, 5, 0},
		{"y is zero", 5, 0, 0},
		{"both zero", 0, 0, 0},
		{"large numbers", math.MaxUint64, 2, float64(math.MaxUint64) / 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Divide2f64(tt.x, tt.y))
		})
	}
}

func TestF64WithDigits(t *testing.T) {
	tests := []struct {
		name     string
		v        float64
		dig      int
		expected float64
	}{
		{"round to 2 digits", 3.14159, 2, 3.14},
		{"round to 1 digit", 3.14159, 1, 3.1},
		{"round up", 3.15, 1, 3.2},
		{"zero value", 0.0, 2, 0.0},
		{"negative digits", -3.14159, 2, -3.14},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, F64WithDigits(tt.v, tt.dig))
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		name     string
		a, n     uint64
		expected uint64
	}{
		{"zero power", 5, 0, 1},
		{"power one", 5, 1, 5},
		{"normal case", 2, 3, 8},
		{"zero base", 0, 5, 0},
		{"one base", 1, 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Pow(tt.a, tt.n))
		})
	}
}

func TestExp(t *testing.T) {
	tests := []struct {
		name     string
		x        uint64
		expected uint64
	}{
		{"zero", 0, 1},
		{"one", 1, 2},
		{"small number", 2, 7},    // e^2 ≈ 7.389
		{"larger number", 5, 148}, // e^5 ≈ 148.413
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Exp(tt.x))
		})
	}
}

func TestCeilDivide(t *testing.T) {
	tests := []struct {
		name     string
		x, y     uint64
		expected uint64
	}{
		{"exact division", 10, 2, 5},
		{"ceiling case", 11, 2, 6},
		{"x less than y", 2, 3, 1},
		{"x equals y", 5, 5, 1},
		{"zero x", 0, 5, 0},
		{"zero y", 5, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, CeilDivide(tt.x, tt.y))
		})
	}
}

func TestToU32(t *testing.T) {
	tests := []struct {
		name     string
		x        uint64
		expected uint32
	}{
		{"zero", 0, 0},
		{"normal value", 1000, 1000},
		{"max uint32", uint64(math.MaxUint32), math.MaxUint32},
		{"overflow", uint64(math.MaxUint32) + 1, math.MaxUint32},
		{"large overflow", math.MaxUint64, math.MaxUint32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ToU32(tt.x))
		})
	}
}

func TestToU64(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected uint64
		wantErr  bool
	}{
		{"valid number", "12345", 12345, false},
		{"zero", "0", 0, false},
		{"max uint64", "18446744073709551615", math.MaxUint64, false},
		{"negative", "-1", 0, true},
		{"invalid input", "abc", 0, true},
		{"empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToU64(tt.s)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// Benchmark tests
func BenchmarkMax(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Max(uint64(i), uint64(i+1))
	}
}

func BenchmarkMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Min(uint64(i), uint64(i+1))
	}
}

func BenchmarkReduce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Sub(uint64(i+1), uint64(i))
	}
}

func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(uint64(i), uint64(i))
	}
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Random(1000)
	}
}

func BenchmarkDivide2f64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Divide2f64(1000, 10)
	}
}

func BenchmarkF64WithDigits(b *testing.B) {
	for i := 0; i < b.N; i++ {
		F64WithDigits(3.14159, 2)
	}
}

func BenchmarkPow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Pow(2, 10)
	}
}

func BenchmarkExp(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Exp(5)
	}
}

func BenchmarkCeilDivide(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CeilDivide(1000, 3)
	}
}

func BenchmarkToU32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToU32(uint64(i))
	}
}

func BenchmarkToU64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToU64("12345")
	}
}
