package api

// Source represents a source
type Source interface {
	Start()

	MeasurementsChannel() <-chan Measurement
}
