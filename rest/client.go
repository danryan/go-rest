package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client type
type Client struct {
	HTTPClient *http.Client
	Endpoint   *url.URL
	Header     http.Header
	Query      url.Values

	EncoderFunc EncoderFunc
	DecoderFunc DecoderFunc
}

// New returns a Client with a configured endpoint and http.Client
func New(endpoint string, client *http.Client) (*Client, error) {
	e, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return NewFromURL(e, client), nil
}

// NewFromURL returns a Client with a configured endpoint from a url.URL and http.Client
func NewFromURL(endpoint *url.URL, client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	if len(endpoint.Path) > 0 && !strings.HasSuffix(endpoint.Path, "/") {
		endpoint.Path = endpoint.Path + "/"
	}

	return &Client{
		HTTPClient: client,
		Endpoint:   endpoint,
		Header:     make(http.Header),
		Query:      endpoint.Query(),
	}
}

// NewRequest returns a new http.Request object
func (c *Client) NewRequest(meth string, path string, input interface{}) (*http.Request, error) {
	ref, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.Endpoint.ResolveReference(ref)

	buf := new(bytes.Buffer)

	if input != nil {
		if err := c.Encode(input, buf); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(meth, u.String(), buf)
	if err != nil {
		return nil, err
	}

	for k := range c.Header {
		req.Header.Set(k, c.Header.Get(k))
	}

	return req, nil
}

// Do performs the request.
func (c *Client) Do(req *http.Request, output interface{}) (*http.Response, error) {

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// fmt.Printf("%v %v - %v\n", req.Method, req.URL, res.StatusCode) // TODO: remove this

	if isHTTPError(res.StatusCode) {
		err = &HTTPError{Response: res}
		return res, err
	}

	if output != nil {
		if w, ok := output.(io.Writer); ok {
			io.Copy(w, res.Body)
		} else {
			err = c.Decode(output, res.Body)
		}
	}

	return res, err
}

func (c *Client) Get(path string, output interface{}) (*http.Response, error) {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, output)
}

// Post makes a POST request
func (c *Client) Post(path string, input, output interface{}) (*http.Response, error) {
	req, err := c.NewRequest("POST", path, input)
	if err != nil {
		return nil, err
	}

	return c.Do(req, output)
}

// Put makes a PUT request
func (c *Client) Put(path string, input, output interface{}) (*http.Response, error) {
	req, err := c.NewRequest("PUT", path, input)
	if err != nil {
		return nil, err
	}

	return c.Do(req, output)
}

// Delete makes a DELETE request
func (c *Client) Delete(path string, input, output interface{}) (*http.Response, error) {
	req, err := c.NewRequest("DELETE", path, input)
	if err != nil {
		return nil, err
	}

	return c.Do(req, output)
}

// HTTPError represents an API error.
type HTTPError struct {
	Response *http.Response
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("%v %v: %d", e.Response.Request.Method, e.Response.Request.URL, e.Response.StatusCode)
}

// isHTTPError determines whether the response status code should be considered an error.
func isHTTPError(status int) bool {
	switch {
	case status > 199 && status < 300:
		return false
	case status == 304:
		return false
	case status == 0:
		return false
	}
	return true
}
