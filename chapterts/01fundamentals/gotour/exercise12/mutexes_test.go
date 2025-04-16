package exercise12

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"sync"
	"testing"
	"time"
)

func TestCrawl_DepthZero(t *testing.T) {
	fetcher := FakeFetcher{
		"http://example.com": &FakeResult{"body", []string{}},
	}
	guard := NewDoubleGuard()
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("http://example.com", 0, fetcher, &guard, &wg)
	wg.Wait()

	if guard.Get("http://example.com") {
		t.Errorf("Crawl with depth 0 should not visit any URL, but visited http://example.com")
	}
}

func TestCrawl_SinglePage(t *testing.T) {
	fetcher := FakeFetcher{
		"http://example.com": &FakeResult{"body", []string{}},
	}
	guard := NewDoubleGuard()
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("http://example.com", 1, fetcher, &guard, &wg)
	wg.Wait()

	if !guard.Get("http://example.com") {
		t.Errorf("Crawl with depth 1 should visit the starting URL, but did not visit http://example.com")
	}
}

func TestCrawl_MultiplePages(t *testing.T) {
	fetcher := FakeFetcher{
		"http://example.com":       &FakeResult{"body1", []string{"http://example.com/page1", "http://example.com/page2"}},
		"http://example.com/page1": &FakeResult{"body2", []string{}},
		"http://example.com/page2": &FakeResult{"body3", []string{}},
	}
	guard := NewDoubleGuard()
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("http://example.com", 2, fetcher, &guard, &wg)
	wg.Wait()

	expected := []string{"http://example.com", "http://example.com/page1", "http://example.com/page2"}
	for _, url := range expected {
		if !guard.Get(url) {
			t.Errorf("Crawl with depth 2 should visit %s, but did not", url)
		}
	}
}

func TestCrawl_CircularLinks(t *testing.T) {
	fetcher := FakeFetcher{
		"http://example.com":       &FakeResult{"body1", []string{"http://example.com/page1"}},
		"http://example.com/page1": &FakeResult{"body2", []string{"http://example.com"}},
	}
	guard := NewDoubleGuard()
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("http://example.com", 2, fetcher, &guard, &wg)
	wg.Wait()

	if !guard.Get("http://example.com") {
		t.Errorf("Crawl should visit http://example.com")
	}
	if !guard.Get("http://example.com/page1") {
		t.Errorf("Crawl should visit http://example.com/page1")
	}
	// Should not infinitely loop due to the guard
}

func TestCrawl_ErrorFetching(t *testing.T) {
	fetcher := FakeFetcher{
		"http://example.com": &FakeResult{"body1", []string{"http://error.com"}},
	}
	// Simulate a fetch error
	fetcher["http://error.com"] = nil

	guard := NewDoubleGuard()
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl("http://example.com", 2, fetcher, &guard, &wg)
	wg.Wait()

	if !guard.Get("http://example.com") {
		t.Errorf("Crawl should visit the initial URL even if subsequent fetches fail")
	}
	if guard.Get("http://error.com") {
		t.Errorf("Crawl should not mark a URL as visited if fetching fails")
	}
}

func TestCrawl_Concurrency(t *testing.T) {
	// Create a test server to simulate fetching with delays
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond) // Simulate network delay
		if r.URL.Path == "/" {
			fmt.Fprintln(w, "Root page", `<a href="/page1">Page 1</a>`, `<a href="/page2">Page 2</a>`)
		} else if r.URL.Path == "/page1" {
			fmt.Fprintln(w, "Page 1", `<a href="/">Root</a>`)
		} else if r.URL.Path == "/page2" {
			fmt.Fprintln(w, "Page 2")
		} else {
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	baseURL := server.URL

	// Create a real Fetcher using http
	realFetcher := &HTTPFetcher{}
	guard := NewDoubleGuard()
	var wg sync.WaitGroup
	wg.Add(1)
	Crawl(baseURL, 2, realFetcher, &guard, &wg)
	wg.Wait()

	expected := []string{baseURL, baseURL + "/page1", baseURL + "/page2"}
	for _, url := range expected {
		if !guard.Get(url) {
			t.Errorf("Crawl should visit %s", url)
		}
	}
}

// HTTPFetcher is a simple Fetcher that uses the net/http package.
type HTTPFetcher struct{}

func (f *HTTPFetcher) Fetch(u string) (body string, urls []string, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bodyBytes := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			bodyBytes = append(bodyBytes, buf[:n]...)
		}
		if err != nil {
			break
		}
	}
	body = string(bodyBytes)

	// Use regular expression for more robust URL extraction
	re := regexp.MustCompile(`<a href="([^"]+)"`)
	matches := re.FindAllStringSubmatch(body, -1)
	for _, match := range matches {
		if len(match) > 1 {
			parsedURL, err := url.Parse(match[1]) // Use url.Parse
			if err == nil {
				absoluteURL := resolveURL(u, parsedURL)
				urls = append(urls, absoluteURL)
			}
		}
	}

	return body, urls, nil
}

func resolveURL(baseURL string, parsedURL *url.URL) string {
	base, err := url.Parse(baseURL)
	if err != nil {
		return parsedURL.String()
	}
	return base.ResolveReference(parsedURL).String()
}
