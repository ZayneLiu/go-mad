package hermod

// AddRequestInterceptor adds a request interceptor.
func (c *Client) AddRequestInterceptor(fn InterceptorFunc[*RequestConfig]) {
	c.requestInterceptors = append(c.requestInterceptors, fn)
}

// AddResponseInterceptor adds a response interceptor.
func (c *Client) AddResponseInterceptor(fn InterceptorFunc[*Response]) {
	c.responseInterceptors = append(c.responseInterceptors, fn)
}
