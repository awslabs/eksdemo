package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Client struct {
	*s3.Client
}

func NewS3Client() *S3Client {
	return &S3Client{s3.NewFromConfig(GetConfig())}
}

func (c *S3Client) CreateBucket(name, region string) error {
	_, err := c.Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(region),
		},
	})

	return err
}

func (c *S3Client) GetBucketLocation(name string) (types.BucketLocationConstraint, error) {
	result, err := c.Client.GetBucketLocation(context.Background(), &s3.GetBucketLocationInput{
		Bucket: aws.String(name),
	})

	if err != nil {
		return types.BucketLocationConstraint(""), err
	}

	return result.LocationConstraint, nil
}

func (c *S3Client) ListBuckets() ([]types.Bucket, error) {
	result, err := c.Client.ListBuckets(context.Background(), &s3.ListBucketsInput{})

	if err != nil {
		return nil, err
	}

	return result.Buckets, nil
}
