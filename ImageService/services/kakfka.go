package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/h2non/bimg"
	"github.com/segmentio/kafka-go"
)

func KConfig() *kafka.Conn {
	topic := os.Getenv("KAFKA_TOPIC")
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" || topic == "" {
		brokers = "localhost:29092"
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
			urls, err := GetImgArray(db, userID, productID)
			if err != nil {
				log.Println(err)
			}
			compImg := WriteToDisk(db, urls)
			if err = AddCompImg(db, userID, productID, compImg); err != nil {
				log.Println(err)
			}
			log.Printf("Compressed images for user %d and product %d saved to DB\n", userID, productID)
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

func WriteToDisk(db *sql.DB, urls []string) (compImages []string) {
	var img Image
	for _, url := range urls {
		if err := img.ImgURLToBuffer(url); err != nil {
			log.Println(err)
		}

		if err := img.ImgCompress(80, bimg.JPEG); err != nil {
			log.Println(err)
		}
		err := img.SaveFile("../compressed")
		if err != nil {
			log.Println(err)
		}
		fPath := filepath.Join("../compressed", img.Name)
		compImages = append(compImages, filepath.Clean(fPath))
		log.Printf("Image %s saved to disk\n", img.Name)
	}
	return compImages
}
