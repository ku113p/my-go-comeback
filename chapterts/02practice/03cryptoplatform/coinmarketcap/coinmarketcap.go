package coinmarketcap

import (
	"crypto/platform/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

const cmcEnvKey = "CMC_API_KEY"
const apiURL = "pro-api.coinmarketcap.com"
const pathLastListing = "/v1/cryptocurrency/listings/latest"
const requestTimeout = 30 * time.Second

type request struct {
	baseURL     string
	path        string
	apiKey      string
	queryParams map[string][]string
}

func newRequest(baseURL, path, apiKey string, queryParams map[string][]string) *request {
	return &request{
		baseURL:     baseURL,
		path:        path,
		apiKey:      apiKey,
		queryParams: queryParams,
	}
}

func newLastListingRequest(apiKey string) *request {
	queryParams := map[string][]string{
		"start":   {"1"},
		"limit":   {"100"},
		"convert": {"USD"},
		"aux":     {"cmc_rank"},
	}
	return newRequest(apiURL, pathLastListing, apiKey, queryParams)
}

func (r *request) url() string {
	v := url.Values(r.queryParams)
	u := url.URL{
		Scheme:  "https",
		Host:    r.baseURL,
		Path:    r.path,
		RawPath: v.Encode(),
	}

	return u.String()
}

func (r *request) fetch() ([]byte, error) {
	url := r.url()
	headers := map[string]string{
		"X-CMC_PRO_API_KEY": r.apiKey,
		"Accept":            "application/json",
	}

	client := &http.Client{Timeout: requestTimeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed request: %s", body)
	}

	return body, nil
}

func GetPrices() ([]*models.TokenPrice, error) {
	apiKey, ok := os.LookupEnv(cmcEnvKey)
	if !ok {
		return nil, fmt.Errorf("env `%s` not found", cmcEnvKey)
	}

	r := newLastListingRequest(apiKey)

	body, err := r.fetch()
	if err != nil {
		return nil, err
	}

	return parsePrices(body)
}

type apiResponse struct {
	Data []struct {
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
		Quote  struct {
			USD struct {
				Price       float64   `json:"price"`
				LastUpdated time.Time `json:"last_updated"`
			} `json:"USD"`
		} `json:"quote"`
	} `json:"data"`
}

func parsePrices(data []byte) ([]*models.TokenPrice, error) {
	var apiResp apiResponse
	if err := json.Unmarshal(data, &apiResp); err != nil {
		return nil, err
	}

	prices := make([]*models.TokenPrice, 0, len(apiResp.Data))
	for _, d := range apiResp.Data {
		tp := models.NewTokenPrice(d.Quote.USD.Price, d.Name, d.Symbol, d.Quote.USD.LastUpdated)
		prices = append(prices, tp)
	}

	return prices, nil
}
