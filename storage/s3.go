package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"product/configs"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Interface interface {
	PresignedPutObject(bucketName string, objectName string) (string, error)
}

type S3 struct {
	client *minio.Client
	bucket string
}

func NewS3() *S3 {

	// Initialize minio client object.
	minioClient, err := minio.New(configs.S3_ENDPOINT, &minio.Options{
		Creds: credentials.NewStaticV4(configs.S3_SECRET_KEY, configs.S3_ACCESS_KEY, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}

	// minio create but

	return &S3{client: minioClient, bucket: configs.S3_BUCKET}
}

func (s *S3) CreateBucketIfNotExists() error {
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

func (s *S3) PresignedPutObject(bucketName string, objectName string) (string, error) {
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

func (s *S3) Put(filename string, file *bytes.Buffer) (string, error) {

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
	return fmt.Sprintf("%s/%s/%s", s.client.EndpointURL(), s.bucket, unique_filename), nil

}
