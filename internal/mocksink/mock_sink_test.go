package mocksink

import (
	"reflect"
	"testing"

	"github.com/renajohn/pac_collector/api"
)

func TestPut(t *testing.T) {
	assertLastValue := func(store MockSink, expect api.Measurement) {
		tail, err := store.LastMeasurement()

		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(tail, expect) {
			t.Errorf("Expected value of %+v got %+v", expect, tail)
		}
	}

	t.Run("Happy case", func(t *testing.T) {
		mockSink := MockSink{}
		measure := api.Measurement{
			MeasurementType: api.WaterTemperature,
			Timestamp:       123456789,
			Value:           []byte("42"),
		}

		mockSink.Put(measure)

		assertLastValue(mockSink, measure)
	})

}
