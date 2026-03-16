package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
}

func New(baseURL string, timeout time.Duration, maxRetries int) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		maxRetries: maxRetries,
	}
}

func (c *Client) Get(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.maxRetries {
				delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				log.Printf("API request failed (attempt %d/%d), retrying in %v...", attempt+1, c.maxRetries+1, delay)
				time.Sleep(delay)
				continue
			}
			break
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
			if attempt < c.maxRetries {
				delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
				log.Printf("API request failed (attempt %d/%d), retrying in %v...", attempt+1, c.maxRetries+1, delay)
				time.Sleep(delay)
				continue
			}
			break
		}

		return body, nil
	}

	if lastErr != nil {
		return nil, fmt.Errorf("request failed after retries: %w", lastErr)
	}
	return nil, fmt.Errorf("request failed after retries")
}

// GetJSON performs a GET and decodes JSON response into dest.
func (c *Client) GetJSON(ctx context.Context, path string, params map[string]string, dest interface{}) error {
	body, err := c.Get(ctx, path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, dest)
}
