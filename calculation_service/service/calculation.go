package service

import (
	"calculation_service/model"
	"math"
)

func CalculateX(data model.RequestData) float64 {
	return round((data.Values[0]/data.Values[1])*data.Values[2], data.E)
}

func CalculateY(data model.RequestData) float64 {
	return round((data.Values[3]/data.Values[4])*data.Values[5], data.E)
}

func round(val float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(val*factor) / factor
}
