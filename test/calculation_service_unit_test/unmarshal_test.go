package calculation_service_unit_test

import (
	"calculation_service/model"
	"testing"
)

func TestValues_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "correct number of values",
			data:    []byte(`[1.0, 2.0, 3.0, 4.0, 5.0, 6.0]`),
			wantErr: false,
		},
		{
			name:    "incorrect number of values",
			data:    []byte(`[1.0, 2.0, 3.0]`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &model.Values{}
			if err := v.UnmarshalJSON(tt.data); (err != nil) != tt.wantErr {
				t.Errorf("Values.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
