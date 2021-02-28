package kafkasink

import (
	"context"
	"fmt"

	"github.com/renajohn/pac_collector/api"
	"github.com/segmentio/kafka-go"
)

type kafkaWriter interface {
	WriteMessages(context context.Context, msgs ...kafka.Message) error
	Close() error
}

type kafkaWriterFactory interface {
	NewWriter() kafkaWriter
}

type kafkaWriterFactoryImpl struct {
	topic    string
	kafkaURL string
}

func (factory *kafkaWriterFactoryImpl) NewWriter() kafkaWriter {
	return &kafka.Writer{
		Addr:     kafka.TCP(factory.kafkaURL),
		Topic:    factory.topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// NewKafkaSink creates a new sink attached to a kafka queue
func NewKafkaSink(sinkURL string, topic string) *KafkaSink {
	factory := kafkaWriterFactoryImpl{
		topic:    topic,
		kafkaURL: sinkURL,
	}
	return newKafkaSinkWithConnectionFactory(&factory)
}

func newKafkaSinkWithConnectionFactory(factory kafkaWriterFactory) *KafkaSink {
	sink := KafkaSink{
		factory: factory,
	}
	return &sink
}

// KafkaSink is saving measurment in a kafka queue
type KafkaSink struct {
	factory kafkaWriterFactory
}

// Put statisfies the api.Sink interface
func (ks *KafkaSink) Put(measurement api.Measurement) error {
	writer := ks.factory.NewWriter()
	defer writer.Close()

	message := kafka.Message{
		Key:   []byte(measurement.MeasurementType),
		Value: measurement.Value,
	}
	writeErr := writer.WriteMessages(context.Background(), message)
	if writeErr != nil {
		fmt.Printf("failed to send a message to Kafka: %g\n", writeErr)
	}

	return writeErr
}
