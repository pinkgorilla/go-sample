package client

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"math"
	"net/http"
	"time"
)

// ClientConfig is http Client configuration
type ClientConfig struct {
	MaxRequestAttempt        int
	MinRequestAttemptDelay   time.Duration
	RequestTimeout           time.Duration
	StopAttemptOnStatusCodes []int
}

// DefaultClientConfig returns default client config
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		MaxRequestAttempt:        3,
		MinRequestAttemptDelay:   200 * time.Millisecond,
		RequestTimeout:           5 * time.Second,
		StopAttemptOnStatusCodes: []int{http.StatusForbidden, http.StatusUnauthorized, http.StatusBadGateway, http.StatusInternalServerError},
	}
}

type Client struct {
	http   *http.Client
	config ClientConfig
}

func NewClient(config *ClientConfig) *Client {
	c := config
	if c == nil {
		dc := DefaultClientConfig()
		c = dc
	}
	hc := &http.Client{
		Timeout: c.RequestTimeout,
	}
	return &Client{
		http:   hc,
		config: *c,
	}
}

// Do executes request
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	copyRequest := c.makeRequestCopier(req)

	var res *http.Response

	attempt := 1
	limit := c.config.MaxRequestAttempt
	if limit < 1 {
		limit = 1
	}

	for {
		r2, err := copyRequest()
		if err != nil {
			return nil, err
		}

		if attempt > limit {
			return res, err
		}
		attempt++

		res, err = c.http.Do(r2)
		if err != nil {
			return res, err
		}

		if res.StatusCode < 200 || res.StatusCode > 300 {
			for _, exStatus := range c.config.StopAttemptOnStatusCodes {
				if res.StatusCode == exStatus {
					return res, err
				}
			}
			d := time.Duration(math.Pow(2, float64(attempt))) * time.Millisecond
			time.Sleep(c.config.MinRequestAttemptDelay + d)
			continue
		}
		return res, err
	}
}

func (c *Client) makeRequestCopier(req *http.Request) func() (*http.Request, error) {
	var bs []byte
	if req.ContentLength > 0 {
		source, err := req.GetBody()
		if err != nil {
			return func() (*http.Request, error) { return nil, err }
		}
		peeker := bufio.NewReader(source)
		bs, _ = peeker.Peek(peeker.Size())
	}
	h := req.Header

	return func() (*http.Request, error) {
		r, e := http.NewRequest(req.Method, req.URL.String(), nil)
		if e != nil {
			return nil, e
		}
		r.Body = ioutil.NopCloser(bytes.NewReader(bs))
		h2 := make(http.Header, len(h))
		for k, vv := range h {
			vv2 := make([]string, len(vv))
			copy(vv2, vv)
			h2[k] = vv2
		}
		r.Header = h2
		return r, nil
	}
}
