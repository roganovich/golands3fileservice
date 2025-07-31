package database

import (
	"log"
	"os"
	"time"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient - глобальная переменная для соединения с S3
var MinioClient *minio.Client

// InitMinio - функция для инициализации соединения с S3
func InitMinio() {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Не удалось открыть соединение с S3: %v", err)
	}

	createBucket()
}

func createBucket() {
	bucketName := os.Getenv("MINIO_ROOT_BUCKET_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := MinioClient.BucketExists(ctx, bucketName)
	if err != nil {
		log.Fatalf("Не удалось создать bucket %v для S3: %v", bucketName, err)
	}

	if !exists {
		err = MinioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Bucket created:", bucketName)
	}
}