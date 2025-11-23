package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

func UploadImageToS3(file multipart.File, fileHeader *multipart.FileHeader, bucket string) (string, error) {
	fileKey := uuid.New().String() + "_" + fileHeader.Filename
	var AWS_REGION = os.Getenv("AWS_REGION")

	contentType := fileHeader.Header.Get("Content-Type")

	_, err := S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &fileKey,
		Body:        file,
		ContentType: &contentType,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, AWS_REGION, fileKey)
	return url, nil
}
