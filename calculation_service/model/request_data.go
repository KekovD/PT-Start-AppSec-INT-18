package model

type Values [6]float64

type RequestData struct {
	Values Values `json:"values"`
	E      int    `json:"e"`
}
