package minio_connection

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKey       string `yaml:"access_key"`
	SecretAccessKey string `yaml:"secret_key"`
	UseSecure       bool   `yaml:"use_secure"`
}

func Connect(cfg *Config) (*minio.Client, error) {

	// Initialize minio client object.
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSecure,
	})
	if err != nil {
		log.Println(err)
	}

	return minioClient, err
}
