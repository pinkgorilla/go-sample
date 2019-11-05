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
	var res *http.Response
	var err error
	var peeker *bufio.Reader
	// cleanUp := func() {
	// 	req.Body.Close()
	// }
	// defer cleanUp()

	attempt := 1
	limit := c.config.MaxRequestAttempt
	if limit < 1 {
		limit = 1
	}

	if req.ContentLength > 0 {
		source, err := req.GetBody()
		if err != nil {
			return nil, err
		}
		peeker = bufio.NewReader(source)
	}

	for {
		if attempt > limit {
			return res, err
		}
		attempt++

		if req.ContentLength > 0 {
			bs, _ := peeker.Peek(peeker.Size())
			req.Body = ioutil.NopCloser(bytes.NewReader(bs))
		}

		res, err = c.http.Do(req)
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
