package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"product/configs"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"google.golang.org/api/option"
)

type S3Interface interface {
	Put(filename string, file *bytes.Buffer) (string, error)
}

type GCS struct {
	client *storage.Client
	bucket string
}

type Minio struct {
	client *minio.Client
	bucket string
}

func NewS3() *GCS {
	ctx := context.Background()

	// Initialize GCS Client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(configs.GCS_CREDENTIALS))
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}

	return &GCS{client: client, bucket: configs.S3_BUCKET}
}

// Upload File to GCS
func (g *GCS) Put(filename string, file *bytes.Buffer) (string, error) {
	ctx := context.Background()
	bucket := g.client.Bucket(g.bucket)

	// Validate file type
	fileType := strings.ToLower(strings.Split(filename, ".")[1])
	if fileType != "jpeg" && fileType != "png" {
		return "", errors.New("file type not allowed")
	}

	// Validate file size (max 10MB)
	if file.Len() > 1024*1024*10 {
		return "", errors.New("file size exceeds the limit (10MB)")
	}

	// Generate a unique filename
	uniqueFilename := fmt.Sprintf("%s-%s.%s", strings.Split(filename, ".")[0], uuid.New().String(), fileType)
	object := bucket.Object(uniqueFilename)

	// Upload file
	writer := object.NewWriter(ctx)
	writer.ContentType = "application/octet-stream"

	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Generate public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucket, uniqueFilename)

	return url, nil
}

func NewMinio() *Minio {
	// Initialize MinIO Client
	minioClient, err := minio.New(configs.S3_ENDPOINT, &minio.Options{
		Creds: credentials.NewStaticV4(configs.S3_ACCESS_KEY, configs.S3_SECRET_KEY, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	return &Minio{client: minioClient, bucket: configs.S3_BUCKET}
}

func (s *Minio) CreateBucketIfNotExists() error {
	exists, err := s.client.BucketExists(context.Background(), s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		err = s.client.MakeBucket(context.Background(), s.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}
	err = s.client.SetBucketPolicy(context.Background(), s.bucket,
		`{
"Version": "2012-10-17",
"Statement": [
{
	"Effect": "Allow",
	"Principal": {
		"AWS": [
			"*"
		]
	},
	"Action": [
		"s3:GetBucketLocation",
		"s3:ListBucket"
	],
	"Resource": [
		"arn:aws:s3:::`+s.bucket+`"
	]
},
{
	"Effect": "Allow",
	"Principal": {
		"AWS": [
			"*"
		]
	},
	"Action": [
		"s3:GetObject"
	],
	"Resource": [
		"arn:aws:s3:::`+s.bucket+`/*"
	]
}
]
}`)

	if err != nil {
		return err
	}
	return nil
}

func (s *Minio) PresignedPutObject(bucketName string, objectName string) (string, error) {
	// Ensure the bucket exists
	err := s.CreateBucketIfNotExists()
	if err != nil {
		return "", err
	}
	// Generate presigned URL
	presignedURL, err := s.client.PresignedPutObject(context.Background(), bucketName, objectName, 24*time.Hour)
	if err != nil {
		log.Fatalln(err)
	}

	return presignedURL.String(), nil
}

func (s *Minio) Put(filename string, file *bytes.Buffer) (string, error) {

	// ensure bucket exists
	err := s.CreateBucketIfNotExists()
	if err != nil {
		return "", err
	}
	file_type := strings.Split(filename, ".")[1]
	if file_type != "jpeg" && file_type != "png" {
		return "", errors.New("file type not allowed")
	}
	if file.Len() > 1024*1024*10 {
		return "", errors.New("file size exceeds the limit 200mb")
	}
	unique_filename := fmt.Sprintf("%s-%s.%s", strings.Split(filename, ".")[0], uuid.New().String(), file_type)
	_, err = s.client.PutObject(context.Background(), s.bucket, unique_filename, bytes.NewReader(file.Bytes()), int64(file.Len()), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("http://localhost:9000/%s/%s",s.bucket, unique_filename), nil
}
