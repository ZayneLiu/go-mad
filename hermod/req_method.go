package hermod

import "net/http"

// ---- Convenience methods ----

// Get sends a GET request.
func (c *Client) Get(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodGet, URL: url,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Post sends a POST request.
func (c *Client) Post(url string, body any, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodPost, URL: url, Body: body,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Put sends a PUT request.
func (c *Client) Put(url string, body any, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodPut, URL: url, Body: body,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Patch sends a PATCH request.
func (c *Client) Patch(url string, body any, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodPatch, URL: url, Body: body,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Delete sends a DELETE request.
func (c *Client) Delete(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodDelete, URL: url,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Head sends a HEAD request.
func (c *Client) Head(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodHead, URL: url,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Options sends an OPTIONS request.
func (c *Client) Options(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodOptions, URL: url,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Connect sends a CONNECT request.
func (c *Client) Connect(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodConnect, URL: url,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}

// Trace sends a TRACE request.
func (c *Client) Trace(url string, opts ...func(*RequestConfig)) (*Response, error) {
	cfg := RequestConfig{
		Method: http.MethodTrace, URL: url,
		Headers: map[string]string{},
		Params:  map[string]string{},
	}
	for _, o := range opts {
		o(&cfg)
	}
	return c.Request(cfg)
}
