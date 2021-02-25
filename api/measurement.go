package api

// MeasurementType represents the base type for all measurement types
type MeasurementType string

const (
	// WaterTemperature represents the current water temperature of the tank as reported by the system
	WaterTemperature MeasurementType = "waterTemperature"
)

// Measurement holds a given measure at a specific time
type Measurement struct {
	MeasurementType MeasurementType
	Timestamp       int64
	Value           []byte
}
