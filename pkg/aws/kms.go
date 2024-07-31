package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KMSClient struct {
	*kms.Client
}

func NewKMSClient() *KMSClient {
	return &KMSClient{kms.NewFromConfig(GetConfig())}
}

func (c *KMSClient) CreateAlias(aliasName, keyID string) error {
	_, err := c.Client.CreateAlias(context.Background(), &kms.CreateAliasInput{
		AliasName:   aws.String(aliasName),
		TargetKeyId: aws.String(keyID),
	})

	return err
}

func (c *KMSClient) CreateKey() (*types.KeyMetadata, error) {
	result, err := c.Client.CreateKey(context.Background(), &kms.CreateKeyInput{
		KeySpec: types.KeySpecSymmetricDefault,
	})

	if err != nil {
		return nil, err
	}

	return result.KeyMetadata, nil
}

func (c *KMSClient) DescribeKey(keyID string) (*types.KeyMetadata, error) {
	result, err := c.Client.DescribeKey(context.Background(), &kms.DescribeKeyInput{
		KeyId: aws.String(keyID),
	})

	if err != nil {
		return nil, err
	}

	return result.KeyMetadata, nil
}

func (c *KMSClient) ListAliases() ([]types.AliasListEntry, error) {
	keys := []types.AliasListEntry{}
	pageNum := 0

	paginator := kms.NewListAliasesPaginator(c.Client, &kms.ListAliasesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		keys = append(keys, out.Aliases...)
		pageNum++
	}

	return keys, nil
}

func (c *KMSClient) ListKeys() ([]types.KeyListEntry, error) {
	keys := []types.KeyListEntry{}
	pageNum := 0

	paginator := kms.NewListKeysPaginator(c.Client, &kms.ListKeysInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		keys = append(keys, out.Keys...)
		pageNum++
	}

	return keys, nil
}
