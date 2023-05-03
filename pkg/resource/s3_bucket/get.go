package s3_bucket

import (
	"errors"
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

type Getter struct {
	s3Client *aws.S3Client
}

func NewGetter(organizationsClient *aws.S3Client) *Getter {
	return &Getter{organizationsClient}
}

func (g *Getter) Init() {
	if g.s3Client == nil {
		g.s3Client = aws.NewS3Client()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var bucket types.Bucket
	var buckets []types.Bucket
	var err error

	if name != "" {
		bucket, err = g.GetBucketByName(name)
		buckets = []types.Bucket{bucket}
	} else {
		buckets, err = g.s3Client.ListBuckets()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(buckets))
}

func (g *Getter) GetBucketByName(name string) (types.Bucket, error) {
	_, err := g.s3Client.GetBucketLocation(name)

	if err != nil {
		// aws-sdk-go-v2 returns a generic APIError instead of defined errors such as types.NoSuchBucket
		var apiErr smithy.APIError
		var nsb types.NoSuchBucket
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == nsb.ErrorCode() {
			return types.Bucket{}, resource.NotFoundError(fmt.Sprintf("bucket %q does not exist", name))
		}

		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "AccessDenied" {
			return types.Bucket{}, fmt.Errorf("access denied, bucket %q not owned by you", name)
		}

		return types.Bucket{}, err
	}

	buckets, err := g.s3Client.ListBuckets()
	if err != nil {
		return types.Bucket{}, err
	}

	for _, b := range buckets {
		if awssdk.ToString(b.Name) == name {
			return b, nil
		}
	}

	return types.Bucket{}, fmt.Errorf("bucket %q not found", name)
}
