package store

// MeasurementType represents the base type for all measurement types
type MeasurementType string

const (
	// WaterTemperature represents the current water temperature of the tank as reported by the system
	WaterTemperature MeasurementType = "waterTemperature"
)

// Measurement holds a given measure at a specific time
type Measurement struct {
	MeasurementType MeasurementType
	TimestampSecond int64
	Measure         float64
}

// Store represents the sink for measurements
type Store interface {
	Put(m Measurement) error
}
