package mockstore

import (
	"reflect"
	"testing"

	"github.com/renajohn/pac_collector/api"
)

func TestPut(t *testing.T) {
	assertLastValue := func(store MockStore, expect api.Measurement) {
		tail, err := store.LastMeasurement()

		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(tail, expect) {
			t.Errorf("Expected value of %+v got %+v", expect, tail)
		}
	}

	t.Run("Happy case", func(t *testing.T) {
		mockStore := MockStore{}
		measure := api.Measurement{
			MeasurementType: api.WaterTemperature,
			TimestampSecond: 123456789,
			Measure:         42.0,
		}

		mockStore.Put(measure)

		assertLastValue(mockStore, measure)
	})

}
