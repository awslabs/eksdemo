package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"
)

type EFSClient struct {
	*efs.Client
}

func NewEFSClient() *EFSClient {
	return &EFSClient{efs.NewFromConfig(GetConfig())}
}

func (c *EFSClient) DescribeFileSystems(fileSystemID string) ([]types.FileSystemDescription, error) {
	input := efs.DescribeFileSystemsInput{}

	if fileSystemID != "" {
		input.FileSystemId = aws.String(fileSystemID)
	}

	result, err := c.Client.DescribeFileSystems(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return result.FileSystems, nil
}
