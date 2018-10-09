package event

// Event ...
type Event struct {
	Source string      `json:"source"`
	Name   string      `json:"name"`
	Data   interface{} `json:"data"`
}
