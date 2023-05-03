package cloudformation_stack

import (
	"errors"
	"fmt"
	"os"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/smithy-go"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Getter struct {
	cloudformationClient *aws.CloudformationClient
}

func NewGetter(cloudformationClient *aws.CloudformationClient) *Getter {
	return &Getter{cloudformationClient}
}

func (g *Getter) Init() {
	if g.cloudformationClient == nil {
		g.cloudformationClient = aws.NewCloudformationClient()
	}
}

func (g *Getter) Get(stackName string, output printer.Output, options resource.Options) error {
	var err error
	var stacks []types.Stack
	clusterName := options.Common().ClusterName

	if clusterName != "" {
		stacks, err = g.GetStacksByCluster(clusterName, stackName)
	} else {
		stacks, err = g.GetStacks(stackName)
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(stacks))
}

func (g *Getter) GetStacks(stackName string) ([]types.Stack, error) {
	stacks, err := g.cloudformationClient.DescribeStacks(stackName)

	var ae smithy.APIError
	if err != nil && errors.As(err, &ae) && ae.ErrorCode() == "ValidationError" {
		return nil, resource.NotFoundError(fmt.Sprintf("cloudformation stack %q not found", stackName))
	}

	return stacks, err
}

func (g *Getter) GetStacksByCluster(clusterName, stackName string) ([]types.Stack, error) {
	stacks, err := g.cloudformationClient.DescribeStacks(stackName)

	if err != nil || clusterName == "" {
		return stacks, err
	}

	clusterStacks := make([]types.Stack, 0, len(stacks))

	for _, stack := range stacks {
		name := awssdk.ToString(stack.StackName)
		if strings.Contains(name, "eksdemo-"+clusterName+"-") || strings.Contains(name, "eksctl-"+clusterName+"-") {
			clusterStacks = append(clusterStacks, stack)
		}
	}

	return clusterStacks, nil
}
