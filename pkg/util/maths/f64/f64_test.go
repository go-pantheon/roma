package f64

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"positive numbers", 2.0, 1.0, 2.0},
		{"negative numbers", -1.0, -2.0, -1.0},
		{"equal numbers", 1.0, 1.0, 1.0},
		{"with zero", 0.0, 1.0, 1.0},
		{"with NaN", math.NaN(), 1.0, math.NaN()},
		{"with Inf", math.Inf(1), 1.0, math.Inf(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Max(tt.x, tt.y)
			if math.IsNaN(tt.expected) {
				assert.True(t, math.IsNaN(got))
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"positive numbers", 2.0, 1.0, 1.0},
		{"negative numbers", -1.0, -2.0, -2.0},
		{"equal numbers", 1.0, 1.0, 1.0},
		{"with zero", 0.0, 1.0, 0.0},
		{"with NaN", math.NaN(), 1.0, math.NaN()},
		{"both NaN", math.NaN(), math.NaN(), math.NaN()},
		{"with positive Inf", math.Inf(1), 1.0, 1.0},
		{"with negative Inf", math.Inf(-1), 1.0, math.Inf(-1)},
		{"with both Inf", math.Inf(1), math.Inf(-1), math.Inf(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Min(tt.x, tt.y)
			if math.IsNaN(tt.expected) {
				assert.True(t, math.IsNaN(got))
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"positive reduction", 5.0, 3.0, 2.0},
		{"zero result", 3.0, 3.0, 0.0},
		{"negative result", 1.0, 2.0, 0.0},
		{"with zero", 1.0, 0.0, 1.0},
		{"both zero", 0.0, 0.0, 0.0},
		{"with NaN x", math.NaN(), 1.0, math.NaN()},
		{"with NaN y", 1.0, math.NaN(), math.NaN()},
		{"with positive Inf", math.Inf(1), 1.0, math.Inf(1)},
		{"with negative Inf", math.Inf(-1), 1.0, 0.0},
		{"both Inf", math.Inf(1), math.Inf(1), 0.0},
		{"tiny numbers", 1e-308, 1e-309, 9e-309},
		{"large numbers", math.MaxFloat64, 1.0, math.MaxFloat64 - 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reduce(tt.x, tt.y)
			if math.IsNaN(tt.expected) {
				assert.True(t, math.IsNaN(got))
			} else {
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float64
		expected float64
	}{
		{"simple addition", 2.0, 3.0, 5.0},
		{"negative numbers", -2.0, -3.0, -5.0},
		{"zero result", 1.0, -1.0, 0.0},
		{"with zero", 1.0, 0.0, 1.0},
		{"both zero", 0.0, 0.0, 0.0},
		{"with NaN x", math.NaN(), 1.0, math.NaN()},
		{"with NaN y", 1.0, math.NaN(), math.NaN()},
		{"positive overflow", math.MaxFloat64, math.MaxFloat64, math.MaxFloat64},
		{"negative overflow", -math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64},
		{"positive + negative Inf", math.Inf(1), math.Inf(-1), math.NaN()},
		{"negative + positive Inf", math.Inf(-1), math.Inf(1), math.NaN()},
		{"both positive Inf", math.Inf(1), math.Inf(1), math.Inf(1)},
		{"both negative Inf", math.Inf(-1), math.Inf(-1), math.Inf(-1)},
		{"tiny numbers", 1e-308, 1e-308, 2e-308},
		{"large numbers below max", math.MaxFloat64 / 2, math.MaxFloat64 / 2, math.MaxFloat64 - 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.x, tt.y)
			if math.IsNaN(got) {
				assert.True(t, math.IsNaN(tt.expected))
			} else if got != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.expected)
			}
		})
	}
}
