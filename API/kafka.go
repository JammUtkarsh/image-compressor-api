package main

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func KConfig() *kafka.Conn {
	topic := os.Getenv("KAFKA_TOPIC")
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" || topic == "" {
		brokers = "localhost:9092"
	}
	conn, err := kafka.DialLeader(context.Background(), "tcp", brokers, topic, 0)
	if err != nil {
		log.Println(err)
	}
	return conn
}

func KProducer(value []byte) {
	conn := KConfig()
	defer conn.Close()
	if _, err := conn.WriteMessages(kafka.Message{Value: value}); err != nil {
		log.Println(err)
	}
}
