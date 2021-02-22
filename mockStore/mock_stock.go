package mockstore

import "github.com/renajohn/pac_collector/store"

// MockStore implements a mock store for testing
type MockStore struct {
	values []store.Measurement
}

// Put a measurement in the sink
func (ms *MockStore) Put(value store.Measurement) error {
	ms.values = append(ms.values, value)
	return nil
}
