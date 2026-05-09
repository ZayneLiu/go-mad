package hermod

import (
	"encoding/base64"
	"time"
)

// Option functions to modify RequestConfig

// WithReqQuery adds query parameters.
func WithReqQuery(params map[string]string) func(*RequestConfig) {
	return func(req *RequestConfig) {
		if req.Params == nil {
			req.Params = make(map[string]string)
		}
		for k, v := range params {
			req.Params[k] = v
		}
	}
}

// WithReqHeader adds a custom header.
func WithReqHeader(key, value string) func(*RequestConfig) {
	return func(req *RequestConfig) {
		if req.Headers == nil {
			req.Headers = make(map[string]string)
		}
		req.Headers[key] = value
	}
}

// WithTimeout sets per‑request timeout.
func WithTimeoutReq(d time.Duration) func(*RequestConfig) {
	return func(req *RequestConfig) {
		req.Timeout = d
	}
}

// WithReqBasicAuth sets HTTP Basic Authentication credentials.
func WithReqBasicAuth(username, password string) func(*RequestConfig) {
	return func(req *RequestConfig) {
		if req.Headers == nil {
			req.Headers = make(map[string]string)
		}
		req.Headers["Authorization"] = basicAuthHeader(username, password)
	}
}

// basicAuthHeader encodes username:password in base64.
func basicAuthHeader(username, password string) string {
	creds := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(creds))
}
