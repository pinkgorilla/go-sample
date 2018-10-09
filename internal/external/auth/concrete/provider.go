package concrete

// Provider ...
type Provider struct {
	host string
}

// Login ...
func (provider *Provider) Login(username, password string) string {
	// code to aquire token from given arguments
	return "my-token"
}

// NewProvider ...
func NewProvider(httpURI string) *Provider {
	return &Provider{
		host: httpURI,
	}
}
