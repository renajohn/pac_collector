package collector

import (
	"reflect"
	"testing"

	"github.com/renajohn/pac_collector/api"
	"github.com/renajohn/pac_collector/internal/mocksink"
)

type MockSource struct {
	measurementsChannel chan api.Measurement
	errorsChannel       chan error
}

func (ms *MockSource) Start() {
}

func (ms *MockSource) MeasurementsChannel() <-chan api.Measurement {
	return ms.measurementsChannel
}

func (ms *MockSource) ErrorsChannel() <-chan error {
	return ms.errorsChannel
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
		source := MockSource{make(chan api.Measurement, 1), make(chan error, 1)}
		mockSink := mocksink.MockSink{}
		measurements := []api.Measurement{{
			MeasurementType: api.WaterTemperature,
			Timestamp:       1,
			Value:           []byte("44"),
		}}

		sendMeasurements(source.measurementsChannel, measurements)

		collector := Collector{Sink: &mockSink, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockSink.Values)
	})

	t.Run("2 values value", func(t *testing.T) {
		source := MockSource{make(chan api.Measurement, 2), make(chan error, 1)}
		mockSink := mocksink.MockSink{}
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

		collector := Collector{Sink: &mockSink, Source: &source}
		collector.Start()

		assertMeasurements(measurements, mockSink.Values)
	})
}
