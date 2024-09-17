package minio

import (
	"bytes"
	"context"

	"github.com/dathuynh1108/clean-arch-base/pkg/config"
	"github.com/dathuynh1108/clean-arch-base/pkg/logger"
	"github.com/dathuynh1108/clean-arch-base/pkg/singleton"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

var (
	store *singleton.Singleton[MinioService]
)

func InitStore() error {
	config := config.GetConfig().MinioConfig
	logger.GetLogger().
		WithFields(logrus.Fields{
			"addr": config.Addr,
			"ssl":  config.SSL,
		}).
		Infof("Init store")
	service, err := NewService(config.Addr, config.AccessKey, config.SecretKey, config.SSL)
	if err != nil {
		return err
	}
	store = singleton.NewSingletonInstance(service)
	return nil
}

func GetStore() MinioService { return store.Get() }

type dao struct {
	client *minio.Client
}

func NewService(addr, accessKey, secretKey string, useSSL bool) (MinioService, error) {
	endpoint := addr
	accessKeyID := accessKey
	secretAccessKey := secretKey
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &dao{client: minioClient}, nil
}

func (d *dao) GetFile(ctx context.Context, bucket, path string) ([]byte, error) {
	obj, err := d.client.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj)
	return buf.Bytes(), err
}

func (d *dao) GetObject(ctx context.Context, bucket, path string) (*minio.Object, error) {
	return d.client.GetObject(ctx, bucket, path, minio.GetObjectOptions{})
}

func (d *dao) PutFile(ctx context.Context, bucket, path string, data []byte, opts minio.PutObjectOptions) (info minio.UploadInfo, err error) {
	file := bytes.NewReader(data)
	info, err = d.client.PutObject(ctx, bucket, path, file, file.Size(), opts)
	if err != nil {
		return
	}
	if info.Location == "" {
		info.Location = d.GetFileURL(ctx, info.Bucket, info.Key)
	}
	return
}

func (d *dao) CreateBucket(ctx context.Context, name string, opts minio.MakeBucketOptions) error {
	err := d.client.MakeBucket(ctx, name, opts)
	if err != nil {
		return err
	}
	return nil
}

func (d *dao) BucketExists(ctx context.Context, name string) bool {
	found, err := d.client.BucketExists(ctx, name)
	if err != nil {
		return false
	}
	return found
}

func (d *dao) RemoveFile(ctx context.Context, bucket, path string, opts minio.RemoveObjectOptions) error {
	return d.client.RemoveObject(ctx, bucket, path, opts)
}

func (d *dao) GetFileURL(ctx context.Context, bucket, path string) string {
	return d.client.EndpointURL().String() + "/" + bucket + "/" + path
}
