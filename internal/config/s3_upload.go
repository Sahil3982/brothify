package config

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

func UploadImageToS3(file multipart.File, fileHeader *multipart.FileHeader, bucket string) (string, error) {
	fileKey := uuid.New().String() + "_" + fileHeader.Filename
	_, err := S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &fileKey,
		Body:        file,
		ContentType: &fileHeader.Header["Content-Type"][0],
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("http://%s.s3.amazonaws.com/%s", bucket, fileKey)

	return url, nil

}
