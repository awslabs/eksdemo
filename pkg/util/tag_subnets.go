package util

import (
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cloudformation_stack"
)

// kubernetes.io/cluster/<clusterName>
const K8stag = `kubernetes.io/cluster/%s`

func GetPrivateSubnets(clusterName string) ([]string, error) {
	stackName := "eksctl-" + clusterName + "-cluster"

	stacks, err := cloudformation_stack.NewGetter(aws.NewCloudformationClient()).GetStacks(stackName)
	if err != nil {
		if _, ok := err.(resource.NotFoundError); ok {
			return nil, fmt.Errorf("cloudformation stack %q not found, is this an eksctl cluster?", stackName)
		}
		return nil, err
	}

	subnets := ""
	for _, o := range stacks[0].Outputs {
		if awssdk.ToString(o.OutputKey) == "SubnetsPrivate" {
			subnets = awssdk.ToString(o.OutputValue)
			continue
		}
	}

	if subnets == "" {
		return nil, fmt.Errorf("no private subnets found in cloudformation stack %q", stackName)
	}

	return strings.Split(subnets, ","), nil
}

func CheckSubnets(clusterName string) error {
	subnets, err := GetPrivateSubnets(clusterName)
	if err != nil {
		return err
	}

	tag := fmt.Sprintf(K8stag, clusterName)
	tagsFilter := []string{tag}

	tags, err := aws.NewEC2Client().DescribeTags(subnets, tagsFilter)
	if err != nil {
		return err
	}

	if len(tags) == 0 {
		return fmt.Errorf("required tag %q not found on any of the following private subnets:\n%s", tag, strings.Join(subnets, "\n"))
	}

	return nil
}

func TagSubnets(clusterName string) error {
	subnets, err := GetPrivateSubnets(clusterName)
	if err != nil {
		return err
	}

	tags := map[string]string{
		fmt.Sprintf(K8stag, clusterName): "",
	}

	fmt.Println("Tagging subnets: " + strings.Join(subnets, ","))
	fmt.Printf("With: %s\n", tags)

	return aws.NewEC2Client().CreateTags(subnets, tags)
}
