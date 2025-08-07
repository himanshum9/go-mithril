package streaming

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Streamer struct {
	writer *kafka.Writer
}

func NewStreamer(brokerAddress string, topic string) *Streamer {
	return &Streamer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{brokerAddress},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (s *Streamer) StreamLocationData(tenantID string, latitude float64, longitude float64, timestamp time.Time) error {
	message := kafka.Message{
		Key:   []byte(tenantID),
		Value: []byte(fmt.Sprintf(`{"latitude": %f, "longitude": %f, "timestamp": "%s"}`, latitude, longitude, timestamp.Format(time.RFC3339))),
	}

	err := s.writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Printf("failed to write message: %v", err)
		return err
	}
	return nil
}

func (s *Streamer) Close() error {
	return s.writer.Close()
}
