package model

type SuccessResponse struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	IsEqual string  `json:"is_equal"`
}
