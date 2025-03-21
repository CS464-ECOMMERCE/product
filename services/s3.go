package services

import (
	"bytes"
	"product/storage"
)

type S3Service struct {
}

func NewS3Service() *S3Service {
	return &S3Service{}
}

func (s *S3Service) UploadFile(filename string, file *bytes.Buffer) (string, error) {
	return storage.NewS3().Put(filename, file)
}
