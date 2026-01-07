package config

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"strings"

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

func UploadInvoiceToS3(htmlContent string, resID uuid.UUID) (string, error) {
	log.Println("Uploading invoice to S3 for reservation ID:", resID)
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return "", fmt.Errorf("AWS_S3_BUCKET environment variable not set")
	}
	log.Println("Using S3 bucket:", bucket)

	fileKey := fmt.Sprintf("invoices/reservation_%d.html", resID)
	contentType := "text/html"

	log.Println("Uploading invoice with key:", fileKey)

	// ✅ Correct: Use strings.NewReader to turn string into io.Reader
	_, err := S3.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      &bucket,
		Key:         &fileKey,
		Body:        strings.NewReader(htmlContent),
		ContentType: &contentType,
	})
	if err != nil {
		log.Printf("S3 upload error: %v", err)
		return "", err
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		awsRegion = "us-east-1" // or your default
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, awsRegion, fileKey)
	log.Println("Invoice uploaded to S3 at URL:", url)
	return url, nil
}
