package api

// MeasurementType represents the base type for all measurement types
type MeasurementType string

const (
	// SWCTemperature represents the current temperatures returned by the SWC PAC system
	SWCTemperature MeasurementType = "SWCTemperature"
)

// Measurement holds a given measure at a specific time
type Measurement struct {
	MeasurementType MeasurementType
	Timestamp       int64
	Value           []byte
}
