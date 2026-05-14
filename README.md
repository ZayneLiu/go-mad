# go-mad

A Go library containing useful utilities for HTTP requests and file operations. The project consists of two main packages: **hermod** (HTTP client) and **saga** (file operations).

## Packages

### 1. hermod - HTTP Client

An Axios-style HTTP client for Go that provides a clean, intuitive API for making HTTP requests with support for interceptors, custom headers, and configurable timeouts.

#### Features

- **Axios-style API**: Clean, familiar interface similar to JavaScript's Axios library
- **Request & Response Interceptors**: Hook into requests and responses for logging, authentication, error handling, etc.
- **Base URL Support**: Configure a base URL and use relative paths in requests
- **Default Headers**: Set headers that apply to all requests
- **Per-request Configuration**: Override timeout, headers, and query parameters on a per-request basis
- **Query Parameters**: Easy parameter management for URLs
- **JSON Body Marshaling**: Automatic JSON serialization of request bodies
- **Basic Authentication**: Built-in support for HTTP Basic Auth
- **Custom HTTP Client**: Use your own `*http.Client` for advanced customization

#### Installation

```bash
go get github.com/zayneliu/go-mad
```

#### Quick Start

```go
package main

import (
	"fmt"
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

#### Core Methods

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

#### Interceptors

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

#### Options

Configure the client with these options:

- `WithBaseURL(url string)` - Set base URL for relative paths
- `WithTimeout(d time.Duration)` - Set default request timeout
- `WithHeader(key, value string)` - Add default header for all requests
- `WithBasicAuth(username, password string)` - Set HTTP Basic Auth credentials
- `WithHTTPClient(cl *http.Client)` - Use custom HTTP client

#### Configuration Options

Pass functions to method calls to configure individual requests:

```go
resp, err := client.Post("/users", userData, func(cfg *hermod.RequestConfig) {
	cfg.Headers["Authorization"] = "Bearer token"
	cfg.Params["version"] = "v2"
	cfg.Timeout = 15 * time.Second
})
```

---

### 2. saga - File Operations

A comprehensive file I/O library providing utilities for reading, writing, archiving, and managing files with support for multiple archive formats.

#### Features

- **File Handles**: Simple abstraction over file paths with convenient methods
- **Read/Write Operations**: Support for plain text, JSON, and binary data
- **File Metadata**: Check existence, size, and modification time
- **Copy & Move**: File duplication and relocation
- **Append Operations**: Add content to existing files
- **Archive Support**: 
  - ZIP compression/decompression
  - GZIP compression/decompression
  - TAR archiving/extraction
  - TAR.GZ compression/decompression
- **Hashing**: Generate SHA256 hashes of file contents
- **Security**: Protection against ZIP slip and TAR slip attacks

#### File Handle Operations

**Getting a File Handle:**

```go
file := saga.File("/path/to/file.txt")
```

**Reading Files:**

```go
// Read as bytes
data, err := file.Read()

// Read as string
text, err := file.Text()

// Read as JSON
var data MyStruct
err := file.JSON(&data)

// Get file path
path := file.Path()
```

**Checking File Properties:**

```go
// Check if file exists
exists, err := file.Exists()

// Get file size
size, err := file.Size()

// Get modification time
modTime, err := file.ModTime()

// Get SHA256 hash
hash, err := file.Hash()
```

**Writing Files:**

```go
// Write text
err := file.WriteText("Hello, World!")

// Write JSON
err := file.WriteJSON(map[string]string{"key": "value"})

// Write bytes
err := file.Write([]byte("data"))

// Write from reader
err := file.Write(io.Reader)
```

**File Operations:**

```go
// Append data
err := file.Append([]byte("\nNew line"))

// Copy file
err := file.Copy("/path/to/dest")

// Move file
err := file.Move("/path/to/new/location")

// Delete file
err := file.Delete()

// Get reader
reader, err := file.Reader()

// Get bytes buffer
buffer, err := file.Bytes()
```

#### Archive Operations

**ZIP Operations:**

```go
// Create a ZIP archive
file1 := saga.File("file1.txt")
file2 := saga.File("directory/")
err := saga.Zip("archive.zip", file1, file2)

// Extract a ZIP archive
archive := saga.File("archive.zip")
err := saga.Unzip(archive, "destination/")
```

**GZIP Operations:**

```go
// Compress a single file
file := saga.File("data.txt")
err := saga.Gz("data.txt.gz", file)

// Decompress
archive := saga.File("data.txt.gz")
err := saga.Ungz(archive, "data.txt")
```

**TAR Operations:**

```go
// Create a TAR archive
file1 := saga.File("file1.txt")
file2 := saga.File("directory/")
err := saga.Tar("archive.tar", file1, file2)

// Extract a TAR archive
archive := saga.File("archive.tar")
err := saga.Untar(archive, "destination/")
```

**TAR.GZ Operations:**

```go
// Create compressed TAR archive
file1 := saga.File("file1.txt")
file2 := saga.File("directory/")
err := saga.TarGz("archive.tar.gz", file1, file2)

// Extract compressed TAR archive
archive := saga.File("archive.tar.gz")
err := saga.UntarGz(archive, "destination/")
```

#### Utility Functions

**Path Resolution:**

```go
// Resolve paths with home directory expansion
resolvedPath, err := saga.ResolvePath("~/documents/file.txt")
```

#### Security Features

- **ZIP Slip Protection**: Validates extracted file paths to prevent directory traversal
- **TAR Slip Protection**: Validates extracted file paths in TAR archives

---

## Module Information

- **Module**: `github.com/zayneliu/go-mad`
- **Go Version**: 1.25.9+

## Project Structure

```
go-mad/
├── hermod/
│   ├── client.go          # Main HTTP client implementation
│   ├── interceptor.go     # Request/response interceptors
│   ├── req_method.go      # HTTP method convenience functions
│   ├── req_config.go      # Request configuration
│   └── response.go        # Response type definition
├── saga/
│   ├── file.go            # File handle and operations
│   ├── archive.go         # Archive operations (ZIP, TAR, GZIP)
│   ├── utils.go           # Utility functions
│   ├── file_test.go       # File operations tests
│   └── archive_test.go    # Archive operations tests
├── go.mod                 # Go module definition
└── README.md              # This file
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a pull request.
