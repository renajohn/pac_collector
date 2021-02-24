package mockstore

import (
	"errors"

	"github.com/renajohn/pac_collector/api"
)

// MockStore implements a mock store for testing
type MockStore struct {
	Values []api.Measurement
}

//LastMeasurement returns the last recorded measurement or throws an error if no measurements where recorded
func (ms *MockStore) LastMeasurement() (api.Measurement, error) {
	if len(ms.Values) == 0 {
		return api.Measurement{}, errors.New("Values array is empty")
	}
	return ms.Values[len(ms.Values)-1], nil
}

// Put a measurement in the sink
func (ms *MockStore) Put(value api.Measurement) error {
	ms.Values = append(ms.Values, value)
	return nil
}