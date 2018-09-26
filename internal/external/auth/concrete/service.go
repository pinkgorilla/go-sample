package concrete

type Service struct {
	host string
}

// Login ...
func (service *Service) Login(username, password string) string {
	// code to aquire token from given arguments
	return "my-token"
}

// NewService ...
func NewService(httpURI string) *Service {
	return &Service{
		host: httpURI,
	}
}
