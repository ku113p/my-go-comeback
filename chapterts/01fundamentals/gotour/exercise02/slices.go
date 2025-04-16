package exercise02

func Pic(dx, dy int, pixelFunc func(int, int) uint8) [][]uint8 {
	p := make([][]uint8, dx)

	for i := range p {
		p[i] = make([]uint8, dy)
	}

	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			p[i][j] = pixelFunc(i, j)
		}
	}

	return p
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
