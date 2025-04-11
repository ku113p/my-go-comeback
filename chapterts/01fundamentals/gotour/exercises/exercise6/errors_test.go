package exercise6

import (
	"math"
	"testing"
)

func TestSqrtPositive(t *testing.T) {
	testCases := []struct {
		input    float64
		expected float64
	}{
		{2, math.Sqrt(2)},
		{4, 2},
		{9, 3},
		{16, 4},
		{25, 5},
		{100, 10},
		{0.5, math.Sqrt(0.5)},
		{0, 0},
	}

	for _, tc := range testCases {
		actual, err := Sqrt(tc.input)
		if err != nil {
			t.Errorf("Sqrt(%f) returned an unexpected error: %v", tc.input, err)
		}
		if math.Abs(actual-tc.expected) > 1e-7 {
			t.Errorf("Sqrt(%f) = %f, expected %f", tc.input, actual, tc.expected)
		}
	}
}

func TestSqrtNegative(t *testing.T) {
	negativeInput := -5.0
	_, err := Sqrt(negativeInput)
	if err == nil {
		t.Errorf("Sqrt(%f) should have returned an error for negative input", negativeInput)
	}
	if _, ok := err.(ErrNegativeSqrt); !ok {
		t.Errorf("Sqrt(%f) returned the wrong error type: got %T, expected ErrNegativeSqrt", negativeInput, err)
	}
	expectedErrorMessage := "cannot Sqrt negative number: -5"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Sqrt(%f) returned the wrong error message: got '%s', expected '%s'", negativeInput, err.Error(), expectedErrorMessage)
	}
}

func TestSqrtZero(t *testing.T) {
	input := 0.0
	actual, err := Sqrt(input)
	if err != nil {
		t.Errorf("Sqrt(%f) returned an unexpected error: %v", input, err)
	}
	if actual != 0 {
		t.Errorf("Sqrt(%f) = %f, expected %f", input, actual, 0.0)
	}
}
