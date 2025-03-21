package storage

import (
	"context"
	"log"
	"product/configs"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Interface interface {
	PresignedPutObject(bucketName string, objectName string) (string, error)
}

type S3 struct {
	client *minio.Client
}

func NewS3() *S3 {

	// Initialize minio client object.
	minioClient, err := minio.New(configs.S3_ENDPOINT, &minio.Options{
		Creds: credentials.NewStaticV4(configs.S3_SECRET_KEY, configs.S3_ACCESS_KEY, ""),
	})
	if err != nil {
		log.Fatalln(err)
	}
	return &S3{client: minioClient}
}

func (s *S3) PresignedPutObject(bucketName string, objectName string) (string, error) {
	// Ensure the bucket exists
	exists, err := s.client.BucketExists(context.Background(), bucketName)
	if err != nil {
		return "", err
	}
	if !exists {
		err = s.client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}

		err = s.client.SetBucketPolicy(context.Background(), bucketName,
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
                "arn:aws:s3:::`+bucketName+`"
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
                "arn:aws:s3:::`+bucketName+`/*"
            ]
        }
    ]
}`)
	}
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
