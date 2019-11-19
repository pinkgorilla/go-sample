package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func HTTPHandlerString(w http.ResponseWriter, r *http.Request) {
	log.Println("string")
	service := NewService()
	str, err := service.String()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write([]byte(str))
}

func HTTPHandlerJSON(w http.ResponseWriter, r *http.Request) {
	log.Println("data")
	service := NewService()
	data, err := service.Data()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(data)
}
