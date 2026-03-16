package service

import (
	"math"
	"testing"
)

func TestFormatCOP(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0, "$ 0"},
		{100, "$ 100"},
		{1000, "$ 1.000"},
		{15000000, "$ 15.000.000"},
		{4200000, "$ 4.200.000"},
		{980000, "$ 980.000"},
		{450, "$ 450"},
		{1234567, "$ 1.234.567"},
		{-5000000, "-$ 5.000.000"},
		{800.5, "$ 801"},     // Rounds to integer
		{999.4, "$ 999"},     // Rounds down
		{999.5, "$ 1.000"},   // Rounds up
		{math.Inf(1), "$ 0"}, // Infinity handled
		{math.NaN(), "$ 0"},  // NaN handled
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatCOP(tt.input)
			if result != tt.expected {
				t.Errorf("formatCOP(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
