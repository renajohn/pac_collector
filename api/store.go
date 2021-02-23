package api

// Store represents the sink for measurements
type Store interface {
	Put(m Measurement) error
}
