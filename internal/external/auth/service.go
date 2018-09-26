package auth

// Service ...
type Service interface {
	Login(username, password string) string
}
