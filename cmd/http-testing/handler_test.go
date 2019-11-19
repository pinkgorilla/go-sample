package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_StringHandler(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(HTTPHandlerJSON))
	c := s.Client()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("expected status ok")
	}
}

func Test_DataHandler(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(HTTPHandlerJSON))
	c := s.Client()
	req, err := http.NewRequest("GET", s.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal("expected status ok")
	}
}
