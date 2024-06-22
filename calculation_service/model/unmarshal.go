package model

import (
	"encoding/json"
	"fmt"
)

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
