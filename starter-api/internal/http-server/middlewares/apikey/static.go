package apikey

import "github.com/pinkgorilla/go-sample/starter-api/errors"

// StaticAPIKeyValidator implementation of
type StaticAPIKeyValidator struct {
}

// Validate validates key against a store
func (m StaticAPIKeyValidator) Validate(key string) (*Identity, error) {
	if key == "" {
		return nil, errors.NewAuthError("unauthorized")
	}
	return &Identity{
		Key:  key,
		Name: "shared-secret",
	}, nil
}
