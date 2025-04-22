package main

import (
	"notes/api/server"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()

		s := server.NewServer([4]byte{127, 0, 0, 1}, 8090)
		s.Run()
	}()
}
