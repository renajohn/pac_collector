package kafkasink

import (
	"context"
	"errors"
	"testing"

	"github.com/renajohn/pac_collector/api"
	"github.com/segmentio/kafka-go"
)

type mockWriter struct {
	messages    []kafka.Message
	returnError bool
}

func (mc *mockWriter) WriteMessages(context context.Context, messages ...kafka.Message) error {
	if mc.returnError {
		return errors.New("something went wrong")
	}

	for _, message := range messages {
		mc.messages = append(mc.messages, message)
	}
	return nil
}

func (mc *mockWriter) Close() error {
	return nil
}

type mockWriterFactoryImpl struct {
	writer      *mockWriter
	returnError bool
	count       int
}

func (kf *mockWriterFactoryImpl) NewWriter() kafkaWriter {
	kf.count++

	writer := mockWriter{
		returnError: kf.returnError,
	}
	kf.writer = &writer
	return &writer
}

func TestPut(t *testing.T) {
	t.Run("Happy case", func(t *testing.T) {
		factory := mockWriterFactoryImpl{}
		sink := newKafkaSinkWithConnectionFactory(&factory)
		measure := api.Measurement{
			MeasurementType: api.WaterTemperature,
			Timestamp:       123456789,
			Value:           []byte("42"),
		}

		sink.Put(measure)

		if factory.count != 1 {
			t.Errorf("Expected 1 connection, got %d", factory.count)
		}
		if string(factory.writer.messages[0].Value) != string(measure.Value) {
			t.Errorf("Expected value of %s, got %s", string(measure.Value), string(factory.writer.messages[0].Value))
		}
	})

	t.Run("If connection returns an error, propagate error", func(t *testing.T) {
		factory := mockWriterFactoryImpl{
			returnError: true,
		}
		sink := newKafkaSinkWithConnectionFactory(&factory)
		measure := api.Measurement{
			MeasurementType: api.WaterTemperature,
			Timestamp:       123456789,
			Value:           []byte("42"),
		}

		err := sink.Put(measure)

		if factory.count != 1 {
			t.Errorf("Expected 1 connection, got %d", factory.count)
		}
		if err == nil {
			t.Errorf("Expected sink.Put to return an error, got nil")
		}
	})

	t.Run("When NewKafkaSink is used, topic and sink URL is persisted", func(t *testing.T) {
		sink := NewKafkaSink("http://tests", "test_topic")

		kafkaURL := sink.factory.(*kafkaWriterFactoryImpl).kafkaURL
		if kafkaURL != "http://tests" {
			t.Errorf("Expected kafka URL of \"http://tests\" and got %s", kafkaURL)
		}
	})

}
