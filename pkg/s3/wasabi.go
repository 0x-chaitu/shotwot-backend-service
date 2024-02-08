package s3

import (
	"context"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

type S3Client struct {
	client *s3.Client
	bucket string
}

func NewWasabiBucket(uri string) (*S3Client, error) {
	sdkConfig, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("eu-central-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("CYH4MC8ZAIQCNKXS0I9C", "T5ZBEVRccXndzP0wgkGPm2L1VcLHeeJVZr9WvJLD", "TOKEN")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(sdkConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("https://s3.eu-central-1.wasabisys.com/")
	})
	return &S3Client{
		client: client,
		bucket: "shotwotwasabitest",
	}, nil

}

func (client *S3Client) PresignedUrl(key string, fileType string) (*v4.PresignedHTTPRequest, error) {
	x := s3.NewPresignClient(client.client)

	request, err := x.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:          aws.String(client.bucket),
		Key:             aws.String(key),
		ContentType:     aws.String(fileType),
		ContentEncoding: aws.String("base64"),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(1000 * int64(time.Second))
	})
	if err != nil {
		return nil, err
	}
	return request, nil
}
