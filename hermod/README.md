# hermod - HTTP Client

An HTTP client for Go that provides a clean, intuitive API for making HTTP requests with support for interceptors, custom headers, and configurable timeouts.

## Features

- **Intuitive API**: Clean, familiar interface similar to JavaScript's Axios library
- **Request & Response Interceptors**: Hook into requests and responses for logging, authentication, error handling, etc.
- **Base URL Support**: Configure a base URL and use relative paths in requests
- **Default Headers**: Set headers that apply to all requests
- **Per-request Configuration**: Override timeout, headers, and query parameters on a per-request basis
- **Query Parameters**: Easy parameter management for URLs
- **JSON Body Marshaling**: Automatic JSON serialization of request bodies
- **Basic Authentication**: Built-in support for HTTP Basic Auth
- **Custom HTTP Client**: Use your own `*http.Client` for advanced customization

## Installation

```bash
go get github.com/zayneliu/go-mad
```

## Quick Start

```go
package main

import (
	"fmt"
	"time"
	"github.com/zayneliu/go-mad/hermod"
)

func main() {
	// Create a client with default settings
	client := hermod.New(
		hermod.WithBaseURL("https://api.example.com"),
		hermod.WithTimeout(10 * time.Second),
		hermod.WithHeader("User-Agent", "my-app/1.0"),
	)

	// Make a GET request
	resp, err := client.Get("/users", func(cfg *hermod.RequestConfig) {
		cfg.Params["id"] = "123"
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.StatusCode)
	fmt.Printf("Body: %s\n", string(resp.Data))
}
```

## Core Methods

**Creating a Client:**

```go
client := hermod.New(opts ...Option)
```

**HTTP Methods:**

- `Get(url string, opts ...func(*RequestConfig))`
- `Post(url string, body any, opts ...func(*RequestConfig))`
- `Put(url string, body any, opts ...func(*RequestConfig))`
- `Patch(url string, body any, opts ...func(*RequestConfig))`
- `Delete(url string, opts ...func(*RequestConfig))`
- `Head(url string, opts ...func(*RequestConfig))`
- `Options(url string, opts ...func(*RequestConfig))`
- `Connect(url string, opts ...func(*RequestConfig))`
- `Trace(url string, opts ...func(*RequestConfig))`
- `Request(config RequestConfig)` - Generic method for any HTTP request

**Generic Request:**

```go
resp, err := client.Request(hermod.RequestConfig{
	Method: "GET",
	URL: "/data",
	Headers: map[string]string{
		"X-Custom-Header": "value",
	},
	Params: map[string]string{
		"filter": "active",
	},
	Context: ctx,
	Timeout: 5 * time.Second,
})
```

## Interceptors

Interceptors allow you to hook into request/response lifecycle:

```go
// Add request interceptor (e.g., for logging or auth)
client.AddRequestInterceptor(func(cfg *hermod.RequestConfig) error {
	fmt.Printf("Making request to: %s %s\n", cfg.Method, cfg.URL)
	return nil
})

// Add response interceptor (e.g., for error handling)
client.AddResponseInterceptor(func(resp *hermod.Response) error {
	if resp.StatusCode >= 400 {
		return fmt.Errorf("request failed: %s", resp.Status)
	}
	return nil
})
```

## Options

Configure the client with these options:

- `WithBaseURL(url string)` - Set base URL for relative paths
- `WithTimeout(d time.Duration)` - Set default request timeout
- `WithHeader(key, value string)` - Add default header for all requests
- `WithBasicAuth(username, password string)` - Set HTTP Basic Auth credentials
- `WithHTTPClient(cl *http.Client)` - Use custom HTTP client

### Configuration Options

Pass functions to method calls to configure individual requests:

```go
resp, err := client.Post("/users", userData, func(cfg *hermod.RequestConfig) {
	cfg.Headers["Authorization"] = "Bearer token"
	cfg.Params["version"] = "v2"
	cfg.Timeout = 15 * time.Second
})
```
