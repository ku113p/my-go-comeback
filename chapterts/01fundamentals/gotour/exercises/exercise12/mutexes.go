package exercise12

import (
	"fmt"
	"sync"
)

type DoubleGuard struct {
	mu sync.Mutex
	v  map[string]bool
}

func NewDoubleGuard() DoubleGuard {
	return DoubleGuard{v: make(map[string]bool)}
}

func (g *DoubleGuard) Keep(url string) {
	g.mu.Lock()
	g.v[url] = true
	g.mu.Unlock()
}

func (g *DoubleGuard) Get(url string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.v[url]
}

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher, guard *DoubleGuard, twg *sync.WaitGroup) {
	if twg != nil {
		defer twg.Done()
	}

	var wg sync.WaitGroup

	if depth <= 0 || guard.Get(url) {
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	guard.Keep(url)

	fmt.Printf("found: %s %q\n", url, body)

	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, guard, &wg)
	}

	wg.Wait()
}

type FakeFetcher map[string]*FakeResult

type FakeResult struct {
	body string
	urls []string
}

func (f FakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		if res != nil {
			return res.body, res.urls, nil
		} else {
			return "", nil, fmt.Errorf("fetch error: %s - FakeResult is nil", url)
		}
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// var fetcher = FakeFetcher{
// 	"https://golang.org/": &FakeResult{
// 		"The Go Programming Language",
// 		[]string{
// 			"https://golang.org/pkg/",
// 			"https://golang.org/cmd/",
// 		},
// 	},
// 	"https://golang.org/pkg/": &FakeResult{
// 		"Packages",
// 		[]string{
// 			"https://golang.org/",
// 			"https://golang.org/cmd/",
// 			"https://golang.org/pkg/fmt/",
// 			"https://golang.org/pkg/os/",
// 		},
// 	},
// 	"https://golang.org/pkg/fmt/": &FakeResult{
// 		"Package fmt",
// 		[]string{
// 			"https://golang.org/",
// 			"https://golang.org/pkg/",
// 		},
// 	},
// 	"https://golang.org/pkg/os/": &FakeResult{
// 		"Package os",
// 		[]string{
// 			"https://golang.org/",
// 			"https://golang.org/pkg/",
// 		},
// 	},
// }
