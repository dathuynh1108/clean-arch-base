package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

type MinioService interface {
	// GetFile return file bytes (read from buffer)
	GetFile(ctx context.Context, bucket, path string) ([]byte, error)
	// GetObject return object as io.Reader
	GetObject(ctx context.Context, bucket, path string) (*minio.Object, error)
	PutFile(ctx context.Context, bucket, path string, data []byte, opts minio.PutObjectOptions) (info minio.UploadInfo, err error)
	CreateBucket(ctx context.Context, name string, opts minio.MakeBucketOptions) error
	BucketExists(ctx context.Context, name string) bool
	RemoveFile(ctx context.Context, bucket, path string, opts minio.RemoveObjectOptions) error
	GetFileURL(ctx context.Context, bucket, path string) string
}
