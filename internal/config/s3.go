package config

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var S3 *s3.Client

func InitS3() {
	S3 = s3.NewFromConfig(AWSConfig)
}

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

func UploadInvoiceToS3(htmlContent string, resID int) (string, error) {
	bucket := os.Getenv("S3_INVOICE_BUCKET")
	fileKey := fmt.Sprintf("invoices/reservation_%d.html", resID)
	contentType := "text/html"

	_, err := S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &fileKey,
		Body:        os.NewFile(0, htmlContent),
		ContentType: &contentType,
	})
	if err != nil {
		return "", err
	}
	var AWS_REGION = os.Getenv("AWS_REGION")
	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, AWS_REGION, fileKey)
	return url, nil
}
