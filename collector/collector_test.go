package collector

import (
	"reflect"
	"testing"

	"github.com/renajohn/pac_collector/api"
	"github.com/renajohn/pac_collector/internal/mockstore"
)

type MockSource struct {
	measurementsChannel chan api.Measurement
}

func (ms *MockSource) Start() (<-chan api.Measurement, error) {
	return ms.measurementsChannel, nil
}

func sendMeasurements(channel chan api.Measurement, measurements []api.Measurement) {
	for _, measurement := range measurements {
		channel <- measurement
	}
	close(channel)
}

func TestStart(t *testing.T) {
	assertMeasurements := func(expect []api.Measurement, got []api.Measurement) {
		if !reflect.DeepEqual(expect, got) {
			t.Errorf("Expected value of %+v got %+v", expect, got)
		}
	}

	t.Run("Single value", func(t *testing.T) {
		source := MockSource{make(chan api.Measurement, 1)}
		mockStore := mockstore.MockStore{}
		measurements := []api.Measurement{{
			MeasurementType: api.WaterTemperature,
			Timestamp:       1,
			Value:           []byte("44"),
		}}

		sendMeasurements(source.measurementsChannel, measurements)

		collector := Collector{Store: &mockStore, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockStore.Values)
	})

	t.Run("2 values value", func(t *testing.T) {
		source := MockSource{make(chan api.Measurement, 2)}
		mockStore := mockstore.MockStore{}
		measurements := []api.Measurement{{
			MeasurementType: api.WaterTemperature,
			Timestamp:       1,
			Value:           []byte("44"),
		}, {
			MeasurementType: api.WaterTemperature,
			Timestamp:       2,
			Value:           []byte("10"),
		}}

		sendMeasurements(source.measurementsChannel, measurements)

		collector := Collector{Store: &mockStore, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockStore.Values)
	})
}
