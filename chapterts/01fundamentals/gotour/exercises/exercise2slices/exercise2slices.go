package exercise2slices

func Pic(dx, dy int) [][]uint8 {
	p := make([][]uint8, dx)

	for i := range p {
		p[i] = make([]uint8, dy)
	}

	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			p[i][j] = f2(i, j)
		}
	}

	return p
}

func f0(x, y int) uint8 {
	return uint8((x + y) / 2)
}

func f1(x, y int) uint8 {
	return uint8(x * y)
}

func f2(x, y int) uint8 {
	return uint8(x ^ y)
}
