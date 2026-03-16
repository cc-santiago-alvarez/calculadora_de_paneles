package model

import (
	"encoding/json"
	"math"
)

// SafeFloat64 is a float64 that serializes Infinity/NaN as null in JSON.
// This handles the case where Go's math.Inf() can't be serialized by encoding/json.
type SafeFloat64 float64

func (f SafeFloat64) MarshalJSON() ([]byte, error) {
	v := float64(f)
	if math.IsInf(v, 0) || math.IsNaN(v) {
		return []byte("null"), nil
	}
	return json.Marshal(v)
}

func (f *SafeFloat64) UnmarshalJSON(data []byte) error {
	var v *float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v == nil {
		*f = SafeFloat64(math.Inf(1))
	} else {
		*f = SafeFloat64(*v)
	}
	return nil
}

func (f SafeFloat64) Float64() float64 {
	return float64(f)
}
