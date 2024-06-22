package service

import (
	"math"
)

func CalculateValue(value1 float64, value2 float64, value3 float64, e int) float64 {
	return Round((value1/value2)*value3, e)
}

func Round(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}
