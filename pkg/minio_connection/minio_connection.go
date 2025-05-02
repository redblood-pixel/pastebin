package minio_connection

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKey       string
	SecretAccessKey string
	UseSSL          bool
}

func Connect() (*minio.Client, error) {
	endpoint := "localhost:9000"
	accessKeyID := "BPYUnB5lsb2sR8OKZYcL"
	secretAccessKey := "iKScGSBikEX4H4EBfX7PacXTaXfUA7hOmenpVm0S"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Println(err)
	}

	return minioClient, err
}
