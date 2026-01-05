package httpclient

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"time"
)

// Client wraps http.Client with custom settings
type Client struct {
	client *http.Client
}

// NewClient creates a new HTTP client with timeout and optional proxy
func NewClient(timeout time.Duration, proxyURL string) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // For pentesting, skip cert verification
	}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	}

	return &Client{
		client: &http.Client{
			Transport: transport,
			Timeout:   timeout,
		},
	}
}

// Get performs a GET request with context
func (c *Client) Get(rawURL string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", rawURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "MVCTrace/1.0")

	return c.client.Do(req)
}

// GetBody performs GET and returns the response body as string
func (c *Client) GetBody(rawURL string) (string, *http.Response, error) {
	resp, err := c.Get(rawURL)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	body := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			body = append(body, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	return string(body), resp, nil
}
