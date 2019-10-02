package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	client "github.com/pinkgorilla/go-sample/pkg/http"
)

type response struct {
	attempt                     int
	ResponseData                interface{}
	ResponseSuccessOnNthAttempt int
	ResponseSuccessStatusCode   int
	ResponseErrorStatusCode     int
	ResponseDelay               time.Duration
}

func newResponse(attempt int, result interface{}, successCode, errorCode int, delay time.Duration) *response {
	return &response{
		attempt:                     1,
		ResponseSuccessOnNthAttempt: attempt,
		ResponseData:                result,
		ResponseSuccessStatusCode:   successCode,
		ResponseErrorStatusCode:     errorCode,
		ResponseDelay:               delay,
	}
}
func (s *response) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(s.ResponseDelay)
		if s.attempt >= s.ResponseSuccessOnNthAttempt {
			w.WriteHeader(s.ResponseSuccessStatusCode)
			json.NewEncoder(w).Encode(s.ResponseData)
			return
		}
		s.attempt++
		w.WriteHeader(s.ResponseErrorStatusCode)
	}
}
func TestGet_1Attempt_configured0MaxAttempt_ResponseOK_ShouldSuccess(t *testing.T) {
	r := newResponse(1, nil, http.StatusOK, http.StatusBadRequest, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, nil)
	assertNil(t, err)
	cfg := client.DefaultClientConfig()
	cfg.MaxRequestAttempt = 0
	c := client.NewClient(cfg)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusOK)
}
func TestGet_1Attempt_ResponseOK_ShouldSuccess(t *testing.T) {
	r := newResponse(1, nil, http.StatusOK, http.StatusBadRequest, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, nil)
	assertNil(t, err)
	c := client.NewClient(nil)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusOK)
}
func TestGet_1Attempt_ResponseBadRequest_ShouldSuccess(t *testing.T) {
	r := newResponse(5, nil, http.StatusOK, http.StatusBadRequest, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, nil)
	assertNil(t, err)
	c := client.NewClient(nil)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusBadRequest)
}
func TestGet_1Attempt_TimeOut_ShouldError(t *testing.T) {
	r := newResponse(3, nil, http.StatusOK, http.StatusBadRequest, 3*time.Second)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, nil)
	assertNil(t, err)
	cfg := client.DefaultClientConfig()
	cfg.RequestTimeout = 1 * time.Second
	c := client.NewClient(cfg)
	_, err = c.Do(req)
	assertNotNil(t, err)
}

func TestGet_3Attempts_ResponseOK_ShouldSuccess(t *testing.T) {
	r := newResponse(3, nil, http.StatusOK, http.StatusBadRequest, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, nil)
	assertNil(t, err)
	c := client.NewClient(nil)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusOK)
}

func TestGet_StopAttemptWhenUnauthorized_ShouldSuccess(t *testing.T) {
	r := newResponse(3, nil, http.StatusOK, http.StatusUnauthorized, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, nil)
	assertNil(t, err)
	c := client.NewClient(nil)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusUnauthorized)
}

func Test_Post_1Attempt_ResponseOK_ShouldSuccess(t *testing.T) {
	r := newResponse(1, nil, http.StatusOK, http.StatusBadRequest, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, bytes.NewBufferString(`{"data":"hello world"}`))
	assertNil(t, err)

	req.Header.Set("content-type", "application/json")
	c := client.NewClient(nil)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusOK)
}

func TestPost_3Attempts_ResponseOK_ShouldSuccess(t *testing.T) {
	r := newResponse(3, nil, http.StatusOK, http.StatusBadRequest, 0)
	svr := httptest.NewServer(r.handler())
	defer svr.Close()

	req, err := http.NewRequest("GET", svr.URL, bytes.NewBufferString(`{"data":"hello world"}`))
	assertNil(t, err)

	req.Header.Set("content-type", "application/json")
	c := client.NewClient(nil)
	res, err := c.Do(req)
	assertNil(t, err)
	assertEqual(t, res.StatusCode, http.StatusOK)
}

func assertNil(t *testing.T, data interface{}) {
	if data != nil {
		t.Fatalf("expected nil, got: %v", data)
	}
}
func assertNotNil(t *testing.T, data interface{}) {
	if data == nil {
		t.Fatal("expected value, got nil")
	}
}
func assertEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Fatalf("expected %v, got: %v", expected, actual)
	}
}
