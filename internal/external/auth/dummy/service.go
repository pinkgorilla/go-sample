package dummy

type Service struct {
}

// Login ...
func (service *Service) Login(username, password string) string {
	// code to aquire token from given arguments
	return "my-dummy-token"
}

// NewService ...
func NewService() *Service {
	return &Service{}
}
