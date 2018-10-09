package auth

import "errors"

// Service ...
type Service struct {
	provider Provider
}

// NewService ...
func NewService(provider Provider) (*Service, error) {
	if provider == nil {
		return nil, errors.New("auth_provider_required")
	}
	return &Service{
		provider: provider,
	}, nil
}

// Login ...
func (service *Service) Login(username, password string) string {
	return service.provider.Login(username, password)
}
