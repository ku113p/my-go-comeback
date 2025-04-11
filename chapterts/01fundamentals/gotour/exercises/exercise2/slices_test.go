package exercise2

import (
	"reflect"
	"testing"
)

func TestPicWithF0(t *testing.T) {
	dx, dy := 3, 4
	expected := [][]uint8{
		{0, 0, 1, 1},
		{0, 1, 1, 2},
		{1, 1, 2, 2},
	}
	actual := Pic(dx, dy, F0)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Pic(%d, %d) with f0 failed: expected %v, got %v", dx, dy, expected, actual)
	}
}

func TestPicWithF1(t *testing.T) {
	dx, dy := 2, 3
	expected := [][]uint8{
		{0, 0, 0},
		{0, 1, 2},
	}
	actual := Pic(dx, dy, F1)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Pic(%d, %d) with f1 failed: expected %v, got %v", dx, dy, expected, actual)
	}
}

func TestPicWithF2(t *testing.T) {
	dx, dy := 4, 2
	expected := [][]uint8{
		{0, 1},
		{1, 0},
		{2, 3},
		{3, 2},
	}
	actual := Pic(dx, dy, F2)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Pic(%d, %d) with f2 failed: expected %v, got %v", dx, dy, expected, actual)
	}
}

func TestPicEmptyDimensions(t *testing.T) {
	dx, dy := 0, 5
	expected := [][]uint8{}
	actual := Pic(dx, dy, F2) // Using f2 here is arbitrary, the result should be the same

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Pic(%d, %d) failed for empty dx: expected %v, got %v", dx, dy, expected, actual)
	}

	dx, dy = 5, 0
	expected = make([][]uint8, 5)
	for i := range expected {
		expected[i] = []uint8{}
	}
	actual = Pic(dx, dy, F2) // Using f2 here is arbitrary

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Pic(%d, %d) failed for empty dy: expected %v, got %v", dx, dy, expected, actual)
	}

	dx, dy = 0, 0
	expected = [][]uint8{}
	actual = Pic(dx, dy, F2) // Using f2 here is arbitrary

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Pic(%d, %d) failed for empty dx and dy: expected %v, got %v", dx, dy, expected, actual)
	}
}

func BenchmarkPic(b *testing.B) {
	dx, dy := 100, 100
	for n := 0; n < b.N; n++ {
		Pic(dx, dy, F2)
	}
}

func BenchmarkPicLarge(b *testing.B) {
	dx, dy := 500, 500
	for n := 0; n < b.N; n++ {
		Pic(dx, dy, F2)
	}
}
