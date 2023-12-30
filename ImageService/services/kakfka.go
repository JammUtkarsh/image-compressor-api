package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/h2non/bimg"
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
	fmt.Println("Kafka Connection Established")
	return conn
}

func KConsumer(db *sql.DB) {
	conn := KConfig()
	defer conn.Close()

	for {
		msg, err := conn.ReadMessage(10e6)
		if err != nil {
			log.Println(err)
			break
		}

		if msg.Value != nil {
			userID, productID := MsgUnarshal(msg.Value)
			WriteToDisk(db, userID, productID)
		}
	}
}

func MsgUnarshal(kMsg []byte) (userID, productID int) {
	msg := struct {
		UserID    int `json:"u"`
		ProductID int `json:"p"`
	}{}

	if err := json.Unmarshal(kMsg, &msg); err != nil {
		return -1, -1
	}

	return msg.UserID, msg.ProductID
}

func WriteToDisk(db *sql.DB, userID, productID int) {
	urls, err := GetImgArray(db, userID, productID)
	if err != nil {
		log.Println(err)
	}
	var img Image
	for _, url := range urls {
		if err := img.ImgURLToBuffer(url); err != nil {
			log.Println(err)
		}

		if err := img.ImgCompress(80, bimg.JPEG); err != nil {
			log.Println(err)
		}
		img.Name = fmt.Sprintf("%d_%d_%s", userID, productID, img.Name)
		err = img.SaveFile("../compressed")
		if err != nil {
			log.Println(err)
		}
		log.Printf("Image %s saved to disk\n", img.Name)
	}
}
