package hermod

import "net/http"

// ---- Convenience methods ----

// Get sends a GET request.
func (c *Client) Get(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{Method: http.MethodGet, URL: url}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Post sends a POST request.
func (c *Client) Post(url string, body any, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{Method: http.MethodPost, URL: url, Body: body}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Put sends a PUT request.
func (c *Client) Put(url string, body any, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{Method: http.MethodPut, URL: url, Body: body}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Patch sends a PATCH request.
func (c *Client) Patch(url string, body any, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{Method: http.MethodPatch, URL: url, Body: body}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Delete sends a DELETE request.
func (c *Client) Delete(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{Method: http.MethodDelete, URL: url}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}
