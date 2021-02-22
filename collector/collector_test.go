package collector

import (
	"reflect"
	"testing"

	"github.com/renajohn/pac_collector/mockstore"
	"github.com/renajohn/pac_collector/store"
)

type MockSource struct {
	measurementsChannel chan store.Measurement
}

func (ms *MockSource) Start() <-chan store.Measurement {
	return ms.measurementsChannel
}

func sendMeasurements(channel chan store.Measurement, measurements []store.Measurement) {
	for _, measurement := range measurements {
		channel <- measurement
	}
	close(channel)
}

func TestStart(t *testing.T) {
	assertMeasurements := func(expect []store.Measurement, got []store.Measurement) {
		if !reflect.DeepEqual(expect, got) {
			t.Errorf("Expected value of %+v got %+v", expect, got)
		}
	}

	t.Run("Single value", func(t *testing.T) {
		source := MockSource{make(chan store.Measurement, 1)}
		mockStore := mockstore.MockStore{}
		measurements := []store.Measurement{{
			MeasurementType: store.WaterTemperature,
			TimestampSecond: 1,
			Measure:         44.0,
		}}

		sendMeasurements(source.measurementsChannel, measurements)

		collector := Collector{Store: &mockStore, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockStore.Values)
	})

	t.Run("2 values value", func(t *testing.T) {
		source := MockSource{make(chan store.Measurement, 2)}
		mockStore := mockstore.MockStore{}
		measurements := []store.Measurement{{
			MeasurementType: store.WaterTemperature,
			TimestampSecond: 1,
			Measure:         44.0,
		}, {
			MeasurementType: store.WaterTemperature,
			TimestampSecond: 2,
			Measure:         10.0,
		}}

		sendMeasurements(source.measurementsChannel, measurements)

		collector := Collector{Store: &mockStore, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockStore.Values)
	})
}
