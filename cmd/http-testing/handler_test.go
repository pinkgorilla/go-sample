package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Response(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := map[string]interface{}{
			"name": "John Doe",
			"age":  10,
		}
		json.NewEncoder(w).Encode(result)

		bs, _ := json.Marshal(result)
		w.Write(bs)

	})

	s := httptest.NewServer(h)
	c := s.Client()

	res, err := c.Get(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	// var y struct {
	// 	Name    string
	// 	Age     int
	// 	Address string
	// }
	json.NewDecoder(res.Body).Decode(&data)
	x := map[string]bool{
		"name":    false, // !is optional
		"age":     false,
		"address": false,
	}
	for key, isOptional := range x {
		_, exists := data[key]
		if !isOptional && !exists {
			t.Fatalf("data not found with key: %s", key)
		}
		// if v != val {
		// 	t.Fatalf("expected value for key %s is %v but get %v", key, val, v)
		// }
	}
	log.Println(data)
}

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
