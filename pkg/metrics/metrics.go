package metrics

type Metrics interface {
	Inc(name string)
}
