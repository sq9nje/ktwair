package ktwair

import (
	"testing"
)

func TestSensorNameToEN(t *testing.T) {
	testCases := []struct {
		name   string
		input  Sensor
		output string
	}{
		{"Pressure", Sensor{Name: "ciśnienie"}, "pressure"},
		{"Temprature", Sensor{Name: "temperatura"}, "temperature"},
		{"Humidity", Sensor{Name: "wilgotność"}, "humidity"},
		{"Unknown string", Sensor{Name: "unknown string"}, "unknown string"},
		{"Empty string", Sensor{Name: ""}, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if r := tc.input.ToEN(); r != tc.output {
				t.Fatalf("Expected translation of '%s' is '%s'. Got '%v'", tc.input.Name, tc.output, r)
			}
		})
	}
}
