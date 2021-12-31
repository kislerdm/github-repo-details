package http

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client http client.
type Client struct {
	c *http.Client
}

// ClientOption defines the option to modify Client's behavior.
type ClientOption func(*Client) error

// ErrorHttp http client error.
type ErrorHttp struct {
	StatusCode int
	Message    string
}

func (e ErrorHttp) Error() string {
	return fmt.Sprintf("[%d] %s", e.StatusCode, e.Message)
}

type Response struct {
	StatusCode int
	Data       []byte
}

// NewClient initiates a Client.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		c: &http.Client{
			Timeout: 0,
		},
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

// WithTimeout sets timeout for the client.
func WithTimeout(t time.Duration) ClientOption {
	if t < 0 {
		return func(c *Client) error {
			return ErrorHttp{
				StatusCode: 504,
				Message:    "set timeout should be positive",
			}
		}
	}
	return func(c *Client) error {
		c.c.Timeout = t
		return nil
	}
}

func (c *Client) setHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		if req.Header.Get(k) == "" {
			req.Header.Add(k, v)
		} else {
			req.Header.Set(k, v)
		}
	}
}

func (c *Client) extractBody(res *http.Response) (*Response, error) {
	if fmt.Sprintf("%d", res.StatusCode)[:1] != "2" {
		return nil, ErrorHttp{
			StatusCode: res.StatusCode,
			Message:    "issue fetching response",
		}
	}
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ErrorHttp{
			StatusCode: 500,
			Message:    err.Error(),
		}
	}
	return &Response{
		StatusCode: res.StatusCode,
		Data:       b,
	}, nil
}

func (c *Client) requestWithContext(ctx context.Context, method string, url string, headers map[string]string, body []byte) (resp *Response, err error) {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, ErrorHttp{req.Response.StatusCode, err.Error()}
	}
	c.setHeaders(req, headers)
	res, err := c.c.Do(req)
	if err != nil {
		return nil, ErrorHttp{res.StatusCode, err.Error()}
	}
	return c.extractBody(res)
}

func (c *Client) request(method string, url string, headers map[string]string, body []byte) (resp *Response, err error) {
	return c.requestWithContext(context.Background(), method, url, headers, body)
}

func (c *Client) GET(url string, headers map[string]string) (resp *Response, err error) {
	return c.request("GET", url, headers, nil)
}

func (c *Client) POST(url string, headers map[string]string, body []byte) (resp *Response, err error) {
	return c.request("POST", url, headers, body)
}

func (c *Client) GETWithContext(ctx context.Context, url string, headers map[string]string) (resp *Response, err error) {
	return c.requestWithContext(ctx, "GET", url, headers, nil)
}

func (c *Client) POSTWithContext(ctx context.Context, url string, headers map[string]string, body []byte) (resp *Response, err error) {
	return c.requestWithContext(ctx, "POST", url, headers, body)
}
