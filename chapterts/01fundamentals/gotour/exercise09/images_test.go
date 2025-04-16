package exercise09

import (
	"image"
	"image/color"
	"reflect"
	"testing"
)

func TestNewImage(t *testing.T) {
	f := func(x, y int) uint8 { return uint8(x + y) }
	img := NewImage(f, 10, 20)

	if img.function == nil {
		t.Errorf("NewImage: function not initialized")
	}
	// Directly comparing functions is not reliable, so we'll test its behavior later
	if img.width != 10 {
		t.Errorf("NewImage: width should be 10, got %d", img.width)
	}
	if img.height != 20 {
		t.Errorf("NewImage: height should be 20, got %d", img.height)
	}
}

func TestImage_ColorModel(t *testing.T) {
	img := NewImage(F0, 10, 10)
	if img.ColorModel() != color.RGBAModel {
		t.Errorf("ColorModel: should be color.RGBAModel, got %v", img.ColorModel())
	}
}

func TestImage_Bounds(t *testing.T) {
	width := 5
	height := 8
	img := NewImage(F0, width, height)
	expectedBounds := image.Rect(0, 0, width, height)
	if img.Bounds() != expectedBounds {
		t.Errorf("Bounds: should be %v, got %v", expectedBounds, img.Bounds())
	}
}

func TestImage_At(t *testing.T) {
	testCases := []struct {
		name     string
		function func(int, int) uint8
		x        int
		y        int
		expected color.Color
	}{
		{
			name:     "F0 test",
			function: F0,
			x:        2,
			y:        4,
			expected: color.RGBA{uint8((2 + 4) / 2), uint8((2 + 4) / 2), 255, 255},
		},
		{
			name:     "F1 test",
			function: F1,
			x:        3,
			y:        5,
			expected: color.RGBA{uint8(3 * 5), uint8(3 * 5), 255, 255},
		},
		{
			name:     "F2 test",
			function: F2,
			x:        6,
			y:        3,
			expected: color.RGBA{uint8(6 ^ 3), uint8(6 ^ 3), 255, 255},
		},
		{
			name:     "Zero coordinates",
			function: F0,
			x:        0,
			y:        0,
			expected: color.RGBA{uint8((0 + 0) / 2), uint8((0 + 0) / 2), 255, 255},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			img := NewImage(tc.function, 10, 10)
			actual := img.At(tc.x, tc.y)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("At(%d, %d): expected %v, got %v", tc.x, tc.y, tc.expected, actual)
			}
		})
	}
}
