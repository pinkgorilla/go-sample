package managed

//Store is an interface providing methods for storing and loading data
type Store interface {
	Push(data interface{}) error
	// Pop removes data from store and returns it to caller
	Pop() (interface{}, error)
	Size() (int, error)
	Dispose()
}
