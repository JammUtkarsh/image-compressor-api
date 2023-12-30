package imageservice

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func KConfig() kafka.ReaderConfig {
	return kafka.ReaderConfig{
		Brokers:  []string{"localhost:29092"},
		GroupID:  "g1",
		Topic:    "original-images",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	}
}

func KConsumer() {
	r := kafka.NewReader(KConfig())
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(string(m.Value))
	}
	r.Close()
}
