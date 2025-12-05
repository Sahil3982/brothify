package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var AWSConfig aws.Config

func InitAWS() {
	var err error

	AWSConfig, err = config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			)),
	)

	if err != nil {
		log.Fatal("❌ Failed to load AWS config:", err)
	}

	log.Println("✔ AWS config loaded")
}
