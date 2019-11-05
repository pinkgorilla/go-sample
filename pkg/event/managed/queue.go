package managed

// Queue is an interface providing methods for pushing and pulling data
type Queue interface {
	Push(data interface{}) error
	Pull() (interface{}, error)
	Dispose()
}
