package minio_storage

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
)

type MinioStorage struct {
	mc         *minio.Client
	bucketName string
}

func NewPastesRepository(mc *minio.Client, bn string) *MinioStorage {
	return &MinioStorage{mc: mc, bucketName: bn}
}

// * Maintain Single responsibility principle

func (r *MinioStorage) CreatePaste(ctx context.Context, name string, expiresAt time.Time, data []byte) error {
	reader := bytes.NewReader(data)
	_, err := r.mc.PutObject(ctx, r.bucketName, name, reader, int64(len(data)), minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"x-amz-expires": expiresAt.Format(time.RFC3339),
		},
	})

	return err
}

func (r *MinioStorage) GetPaste(ctx context.Context, name string) ([]byte, error) {
	object, err := r.mc.GetObject(ctx, r.bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	return io.ReadAll(object)
}

func (r *MinioStorage) DeletePaste(ctx context.Context, name string) error {
	err := r.mc.RemoveObject(ctx, r.bucketName, name, minio.RemoveObjectOptions{})
	return err
}
