package exercise01

import "testing"

func TestSum(t *testing.T) {
	testCases := []struct {
		name          string
		a             float64
		expected_from float64
		expected_to   float64
	}{
		{"Example", 2, 1.414, 1.415},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := Sqrt(tc.a)
			if !(tc.expected_from < actual && actual < tc.expected_to) {
				t.Errorf("Sqrt(%v) returned %v, expected between (%v, %v)", tc.a, actual, tc.expected_from, tc.expected_to)
			}
		})
	}
}
