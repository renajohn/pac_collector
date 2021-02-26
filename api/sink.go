package api

// Sink represents the sink for measurements
type Sink interface {
	Put(m Measurement) error
}
