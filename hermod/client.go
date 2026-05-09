// hermod is an HTTP client that provides an Axios‑style API
package hermod

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// RequestConfig holds all settings for a single HTTP request.
type RequestConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Params  map[string]string // query string
	Body    any               // will be JSON‑marshalled if not nil
	Timeout time.Duration     // per‑request timeout, overrides client's
	Context context.Context
}

// Response wraps the standard *http.Response with helpers.
type Response struct {
	Status     string
	StatusCode int
	Headers    map[string]string
	Data       []byte // raw body
}

// Client is the Axios‑style client.
type Client struct {
	baseURL              string
	defaultHeaders       map[string]string
	timeout              time.Duration
	requestInterceptors  []InterceptorFunc[*RequestConfig]
	responseInterceptors []InterceptorFunc[*Response]
	httpClient           *http.Client
}

// InterceptorFunc is a hook that can modify config/response or return an error.
type InterceptorFunc[T any] func(T) (any, error)

// Option configures a Client.
type Option func(*Client)

// DefaultClient is a ready‑to‑use client with sensible defaults.
var DefaultClient = New()

// New creates a new Axios client with optional base URL, timeout, etc.
func New(opts ...Option) *Client {
	c := &Client{
		defaultHeaders: map[string]string{},
		timeout:        30 * time.Second,
		httpClient:     &http.Client{},
	}
	for _, o := range opts {
		o(c)
	}
	return c
}

// WithBaseURL sets a base URL that will be prepended to relative URLs.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

// WithTimeout sets the default request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		c.timeout = d
	}
}

// WithHeader sets a default header for every request.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.defaultHeaders[key] = value
	}
}

// WithBasicAuth sets HTTP Basic Authentication credentials.
func WithBasicAuth(username, password string) Option {
	return WithHeader("Authorization", basicAuthHeader(username, password))

}

// WithHTTPClient replaces the underlying *http.Client.
func WithHTTPClient(cl *http.Client) Option {
	return func(c *Client) {
		c.httpClient = cl
	}
}

// Request performs a generic HTTP request.
func (c *Client) Request(config RequestConfig) (*Response, error) {
	// 1. Apply request interceptors
	cfg := config
	for _, interceptor := range c.requestInterceptors {
		result, err := interceptor(&cfg)
		if err != nil {
			return nil, err
		}
		var ok bool
		cfg, ok = result.(RequestConfig)
		if !ok {
			return nil, fmt.Errorf("request interceptor must return RequestConfig")
		}
	}

	// 2. Build URL
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}
	if c.baseURL != "" && !u.IsAbs() {
		u, err = url.Parse(c.baseURL + cfg.URL)
		if err != nil {
			return nil, fmt.Errorf("join base URL: %w", err)
		}
	}
	// Query params
	q := u.Query()
	for k, v := range cfg.Params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	// 3. Prepare body
	var bodyReader io.Reader
	if cfg.Body != nil {
		b, err := json.Marshal(cfg.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	// 4. Build request
	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, cfg.Method, u.String(), bodyReader)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	// 5. Headers (defaults + custom)
	for k, v := range c.defaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	// 6. Execute
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = c.timeout
	}
	client := c.httpClient
	if timeout > 0 {
		// create a one‑shot client with the request timeout
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		req = req.WithContext(ctx)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// 7. Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	// 8. Build response
	headers := make(map[string]string)
	for k := range resp.Header {
		headers[k] = resp.Header.Get(k)
	}
	r := &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Data:       body,
	}

	// 9. Apply response interceptors
	for _, interceptor := range c.responseInterceptors {
		result, err := interceptor(r)
		if err != nil {
			return nil, err
		}
		var ok bool
		r, ok = result.(*Response)
		if !ok {
			return nil, fmt.Errorf("response interceptor must return *Response")
		}
	}

	return r, nil
}
