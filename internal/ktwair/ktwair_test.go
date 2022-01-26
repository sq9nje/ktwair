package ktwair

import (
	"testing"
)

func TestSensorNameToEN(t *testing.T) {
	testCases := []struct {
		name   string
		input  string
		output string
	}{
		{"Pressure", "ciśnienie", "pressure"},
		{"Temprature", "temperatura", "temperature"},
		{"Humidity", "wilgotność", "humidity"},
		{"Unknown string", "unknown string", "unknown string"},
		{"Empty string", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := SensorNameToEN(tc.input); r != tc.output {
				t.Fatalf("Expected translation of '%s' is '%s'. Got '%v'", tc.input, tc.output, r)
			}
		})
	}
}
