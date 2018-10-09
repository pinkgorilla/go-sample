package auth

// Provider ...
type Provider interface {
	Login(username, password string) string
}
