package api

// Source represents a source
type Source interface {
	Start() <-chan Measurement
}
