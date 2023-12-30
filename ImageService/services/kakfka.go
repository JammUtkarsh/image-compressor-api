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

func KConfig() kafka.ReaderConfig {
	topic := os.Getenv("KAFKA_TOPIC")
	broker := os.Getenv("KAFKA_BROKERS")
	return kafka.ReaderConfig{
		Brokers:  []string{broker},
		GroupID:  "g1",
		Topic:    topic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	}
}

func KConsumer(db *sql.DB) {
	r := kafka.NewReader(KConfig())
	log.Println("Kafka Reader Initialized")
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			break
		}
		if m.Value != nil {
			userID, productID := MsgUnarshal(m.Value)
			WriteToDisk(db, userID, productID)
		}
	}
	r.Close()
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
		
		if err:=img.ImgCompress(80, bimg.JPEG); err != nil {
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