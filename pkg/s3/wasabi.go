package s3

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func NewWasabiBucket(uri string) (*s3.Client, error) {
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
	return client, nil

	// x := s3.NewPresignClient(client)
	// bucketName := "warehouse-test"
	// objectKey := "warehouse.mp4"
	// request, err := x.PresignPutObject(context.TODO(), &s3.PutObjectInput{
	// 	Bucket:          aws.String(bucketName),
	// 	Key:             aws.String(objectKey),
	// 	ContentType:     aws.String("video/mp4"),
	// 	ContentEncoding: aws.String("base64"),
	// }, func(opts *s3.PresignOptions) {
	// 	opts.Expires = time.Duration(100000 * int64(time.Second))
	// })
	// if err != nil {
	// 	log.Printf("Couldn't get a presigned request to put %v:%v. Here's why: %v\n",
	// 		bucketName, objectKey, err)
	// }
	// logger.Info(request)
}
