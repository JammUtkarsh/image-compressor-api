package services

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
)

const (
	GetImageArrayQuery = "SELECT product_images FROM products WHERE id = $1 AND user_id = $2"
	AddCompImageQuery  = "UPDATE products SET compressed_product_images = $1 WHERE id = $2 AND user_id = $3"
)

var (
	port     = 5432
	host     = os.Getenv("POSTGRES_HOST")
	user     = os.Getenv("POSTGRES_USER")
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = os.Getenv("POSTGRES_DB")
)

func init() {
	if host == "" || user == "" || password == "" || dbname == "" {
		log.Fatalln("Please set the postgres enviroment variables")
	}
}

func Connect() (db *sql.DB, err error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return
	}
	for i := 0; i < 5; i++ {
		if err = db.Ping(); err != nil {
			log.Println("failed to connect to database, retrying...")
			time.Sleep(6 * time.Second)
		} else {
			return
		}
	}
	if err != nil {
		fmt.Println("failed to connect to database")
		os.Exit(1)
	}
	return
}

func GetImgArray(db *sql.DB, id int, userID int) (imageArray []string, err error) {
	err = db.QueryRow(GetImageArrayQuery, id, userID).Scan(pq.Array(&imageArray))
	return
}

func AddCompImg(db *sql.DB, id int, userID int, compImage []string) (err error) {
	_, err = db.Exec(AddCompImageQuery, pq.Array(compImage), id, userID)
	return
}
