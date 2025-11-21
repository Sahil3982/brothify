package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3 *s3.Client

func initS3() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	config.WithRegion(os.Getenv("AWS_REGION"))
	config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""))

	if err != nil {
		log.Fatal("Failed to load AWS config:", err)
	}

	S3 = s3.NewFromConfig(cfg)

}
