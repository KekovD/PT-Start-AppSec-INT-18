package calculation_service_unit_test

import (
	"calculation_service/service"
	"testing"
)

func TestCalculateValue(t *testing.T) {
	tests := []struct {
		name     string
		value1   float64
		value2   float64
		value3   float64
		e        int
		expected float64
	}{
		{
			name:     "Expected 15.0",
			value1:   10.0,
			value2:   2.0,
			value3:   3.0,
			e:        2,
			expected: 15.0,
		},
		{
			name:     "Expected 7.5",
			value1:   5.0,
			value2:   2.0,
			value3:   3.0,
			e:        2,
			expected: 7.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.CalculateValue(tt.value1, tt.value2, tt.value3, tt.e); got != tt.expected {
				t.Errorf("CalculateValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}
