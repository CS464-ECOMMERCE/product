package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"product/configs"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type S3Interface interface {
	Put(filename string, file *bytes.Buffer) (string, error)
}

type GCS struct {
	client *storage.Client
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
