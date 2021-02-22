package mockstore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/renajohn/pac_collector/store"
)

func TestPut(t *testing.T) {
	lastValue := func(values []store.Measurement) (*store.Measurement, error) {
		if len(values) == 0 {
			return nil, errors.New("Values array is empty")
		}
		return &values[len(values)-1], nil
	}
	assertLastValue := func(store MockStore, expect *store.Measurement) {
		tail, err := lastValue(store.values)

		if err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(tail, expect) {
			t.Errorf("Expected value of %+v got %+v", expect, tail)
		}
	}

	t.Run("Happy case", func(t *testing.T) {
		mockStore := MockStore{}
		measure := store.Measurement{
			MeasurementType: store.WaterTemperature,
			TimestampSecond: 123456789,
			Measure:         42.0,
		}

		mockStore.Put(measure)

		assertLastValue(mockStore, &measure)
	})

}
