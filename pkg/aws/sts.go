package aws

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var accountId, partition string

func AccountId() string {
	if accountId == "" {
		getCallerIdentity()
	}
	return accountId
}

func Partition() string {
	if partition == "" {
		getCallerIdentity()
	}
	return partition
}

func getCallerIdentity() {
	out, err := sts.NewFromConfig(GetConfig()).GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get AWS identity: %w", err))
	}

	arn, err := arn.Parse(aws.ToString(out.Arn))
	if err != nil {
		log.Fatal(fmt.Errorf("failed to parse ARN looking up AWS identity: %w", err))
	}

	accountId = arn.AccountID
	partition = arn.Partition
}
