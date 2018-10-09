package business

import "github.com/pinkgorilla/go-sample/internal/external/auth"

// Service ...
type Service struct {
	Actor string
	auth  *auth.Service
}

// NewService ...
func NewService(auth *auth.Service) *Service {
	service := &Service{
		auth: auth,
	}
	return service
}

// GetToken ...
func (service *Service) GetToken(username, password string) string {
	return service.auth.Login(username, password)
}
