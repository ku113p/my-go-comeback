package exercise4

func fibonacci() func() int {
	a, b := -1, 1

	return func() int {
		next := a + b
		a, b = b, next
		return next
	}
}
