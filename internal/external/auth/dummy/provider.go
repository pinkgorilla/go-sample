package dummy

// Provider ...
type Provider struct {
}

// Login ...
func (provider *Provider) Login(username, password string) string {
	// code to aquire token from given arguments
	return "my-dummy-token"
}

// NewProvider ...
func NewProvider() *Provider {
	return &Provider{}
}
