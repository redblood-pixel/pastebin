package minio_storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/google/uuid"
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

func (r *MinioStorage) DeletePastes(ctx context.Context, userID int, pastesID []uuid.UUID) error {

	// S3 massive delete
	objCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objCh)

		for i := range pastesID {

			opts := minio.ListObjectsOptions{
				Prefix:    strconv.Itoa(userID) + "/" + pastesID[i].String(),
				Recursive: true,
			}
			fmt.Println(opts.Prefix)
			for object := range r.mc.ListObjects(ctx, r.bucketName, opts) {
				fmt.Println("key", object.Key, object.Err)
				if object.Err != nil {
					continue
				}
				objCh <- object
			}
		}
	}()

	errorCh := r.mc.RemoveObjects(ctx, r.bucketName, objCh, minio.RemoveObjectsOptions{})
	var err error
	for e := range errorCh {
		err = e.Err
		fmt.Println("err", e)
	}

	return err
}
