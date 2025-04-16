package exercise09

import (
	"image"
	"image/color"
)

type Image struct {
	function func(int, int) uint8
	width    int
	height   int
}

func NewImage(f func(int, int) uint8, w, h int) Image {
	return Image{function: f, width: w, height: h}
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.width, i.height)
}

func (i Image) At(x, y int) color.Color {
	v := i.function(x, y)
	return color.RGBA{v, v, 255, 255}
}

func F0(x, y int) uint8 {
	return uint8((x + y) / 2)
}

func F1(x, y int) uint8 {
	return uint8(x * y)
}

func F2(x, y int) uint8 {
	return uint8(x ^ y)
}
