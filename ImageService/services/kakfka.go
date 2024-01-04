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

// The function establishes a connection to a Kafka broker using the provided topic and brokers, or
// defaults to localhost:29092 if not provided.
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

// The function `KConsumer` reads messages from a Kafka topic, retrieves image URLs from a database
// based on the message content, compresses the images, and saves the compressed images back to the
// database.
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
			urls, err := GetImgArray(db, productID, userID)
			if err != nil {
				log.Println(err)
			}
			compImg := WriteToDisk(urls)
			if err = AddCompImg(db, userID, productID, compImg); err != nil {
				log.Println(err)
			}
			log.Printf("Compressed images for user %d and product %d saved to DB\n", userID, productID)
		}
	}
}

// This takes a byte array representing a JSON message and returns the user ID
// and product ID extracted from the message.
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

// The function takes a database connection and a list of URLs, downloads and compresses the images
// from the URLs, saves them to disk, and returns a list of the file paths of the compressed images.
func WriteToDisk(urls []string) (compImages []string) {
	var img Image
	for _, url := range urls {
		if err := img.ImgURLToBuffer(url); err != nil {
			log.Println(err)
		}

		if err := img.ImgCompress(50, bimg.JPEG); err != nil {
			log.Println(err)
		}
		err := img.SaveFile("../compressed")
		if err != nil {
			log.Println(err)
		}
		pwd, err := os.Getwd()
		if err != nil {
			log.Println(err)
		}
		fPath := filepath.Join(pwd, "../compressed", img.Name)
		compImages = append(compImages, filepath.Clean(fPath))
		log.Printf("Image %s saved to disk\n", img.Name)
	}
	return compImages
}
