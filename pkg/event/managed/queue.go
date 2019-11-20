package managed

// Queue is an interface providing methods for pushing and pulling data
type Queue interface {
	Push(data interface{}) error
	Pull() (interface{}, error)
	Dispose()
}

// Size interface for types with size
type Size interface {
	// Size returns the size of implementing type
	Size() int
}
