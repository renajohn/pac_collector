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

func (ms *MockSource) Start() <-chan api.Measurement {
	return ms.measurementsChannel
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
			TimestampSecond: 1,
			Measure:         44.0,
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
			TimestampSecond: 1,
			Measure:         44.0,
		}, {
			MeasurementType: api.WaterTemperature,
			TimestampSecond: 2,
			Measure:         10.0,
		}}

		sendMeasurements(source.measurementsChannel, measurements)

		collector := Collector{Store: &mockStore, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockStore.Values)
	})
}
