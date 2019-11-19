package main

import "time"

type Service struct {
}
type Data struct {
	Number int    `json:"number"`
	Text   string `json:"text"`
}

// NewService returns new service
func NewService() *Service {
	return &Service{}
}

// String returns string
func (s *Service) String() (string, error) {
	return "Hello", nil
}

// Data returns Data
func (s *Service) Data() (*Data, error) {
	return &Data{
		Number: time.Now().Nanosecond(),
		Text:   "hello",
	}, nil
}
