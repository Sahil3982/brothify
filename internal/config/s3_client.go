package config

import "github.com/aws/aws-sdk-go-v2/service/s3"

var S3 *s3.Client

func InitS3() {
	S3 = s3.NewFromConfig(AWSConfig)
}
