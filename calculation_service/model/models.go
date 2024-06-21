package model

import (
	"encoding/json"
	"fmt"
)

type SuccessResponse struct {
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	IsEqual string  `json:"is_equal"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Values [6]float64

type RequestData struct {
	Values Values `json:"values"`
	E      int    `json:"e"`
}

func (v *Values) UnmarshalJSON(data []byte) error {
	var tmp []float64
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	if len(tmp) != 6 {
		return fmt.Errorf("incorrect number of values: expected 6, got %d", len(tmp))
	}
	for i := range tmp {
		v[i] = tmp[i]
	}
	return nil
}
