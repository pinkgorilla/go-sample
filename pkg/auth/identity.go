package auth

import "context"

// Identity represents the api caller
type Identity struct {
	Type string
	ID   interface{}
	Name string
}

// NewIdentity returns new identity
func NewIdentity(i interface{}, t, n string) *Identity {
	return &Identity{
		ID:   i,
		Type: t,
		Name: n,
	}
}

type k string

const key = k("auth")

// FromContext get app from context
func FromContext(ctx context.Context) *Identity {
	id, ok := ctx.Value(key).(*Identity)
	if !ok {
		return nil
	}
	return id
}

// ToContext put app to context
func ToContext(ctx context.Context, id *Identity) context.Context {
	return context.WithValue(ctx, key, id)
}
