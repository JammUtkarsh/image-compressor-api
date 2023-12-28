package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessID := os.Getenv("MINIO_ROOT_USER")
	accessSecret := os.Getenv("MINIO_ROOT_PASSWORD")
	bName := os.Getenv("MINIO_BUCKET_NAME")
	if endpoint == "" || accessID == "" || accessSecret == "" || bName == "" {
		log.Fatalln("Please set the minIO enviroment variables")
	}
	fmt.Println(accessID)
	fmt.Println(accessSecret)
	minClient, err := GetMinioClient(accessID, accessSecret, endpoint, "", false) // since we are using minio locally, we are not using SSL
	if err != nil {
		log.Fatalln(err)
	}

	if err := CreateBucket(ctx, minClient, minio.MakeBucketOptions{}, bName); err != nil {
		log.Fatalln(err)
	}
	
	files, err := os.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if info, err := UploadFile(ctx, minClient, bName, f.Name(), fmt.Sprintf("./images/%s", f.Name())); err != nil {
			log.Fatalln(err)
		} else {
			log.Printf("Uploaded %s of size %d successfully\n", info.Key, info.Size)
		}
	}
}

func GetMinioClient(accessKeyID, secretAccessKey, endpoint, token string, useSSL bool) (*minio.Client, error) {
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, token),
		Secure: useSSL,
	})
}

func CreateBucket(ctx context.Context, minClient *minio.Client, options minio.MakeBucketOptions, bName string) (err error) {
	exists, errBucketExists := minClient.BucketExists(ctx, bName)
	if errBucketExists == nil && exists {
		log.Printf("%s already exists\n", bName)
		return nil
	}
	if err = minClient.MakeBucket(ctx, bName, options); err != nil {
		return err
	}
	return nil
}

func UploadFile(ctx context.Context, minClient *minio.Client, bName, oName, fName string) (fInfo minio.UploadInfo, err error) {
	contentType := http.DetectContentType([]byte(fName))
	if fInfo, err = minClient.FPutObject(ctx, bName, oName, fName, minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return fInfo, err
	}
	return fInfo, nil
}
