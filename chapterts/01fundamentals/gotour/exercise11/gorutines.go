package exercise11

import "golang.org/x/tour/tree"

func Walk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	if t.Left != nil {
		Walk(t.Left, ch)
	}
	ch <- t.Value
	if t.Right != nil {
		Walk(t.Right, ch)
	}
}

func Same(t1, t2 *tree.Tree) bool {
	t1c, t2c := make(chan int), make(chan int)

	go func() {
		Walk(t1, t1c)
		close(t1c)
	}()
	go func() {
		Walk(t2, t2c)
		close(t2c)
	}()

	for v1 := range t1c {
		if v2, ok := <-t2c; !ok || v1 != v2 {
			return false
		}
	}
	if _, ok := <-t2c; ok {
		return false
	}

	return true
}
