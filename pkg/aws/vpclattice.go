package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice"
	"github.com/aws/aws-sdk-go-v2/service/vpclattice/types"
)

type VPCLatticeClient struct {
	*vpclattice.Client
}

func NewVPCLatticeClient() *VPCLatticeClient {
	return &VPCLatticeClient{vpclattice.NewFromConfig(GetConfig())}
}

func (c *VPCLatticeClient) GetService(id string) (*vpclattice.GetServiceOutput, error) {
	result, err := c.Client.GetService(context.Background(), &vpclattice.GetServiceInput{
		ServiceIdentifier: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VPCLatticeClient) GetServiceNetwork(id string) (*vpclattice.GetServiceNetworkOutput, error) {
	result, err := c.Client.GetServiceNetwork(context.Background(), &vpclattice.GetServiceNetworkInput{
		ServiceNetworkIdentifier: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VPCLatticeClient) GetTargetGroup(id string) (*vpclattice.GetTargetGroupOutput, error) {
	result, err := c.Client.GetTargetGroup(context.Background(), &vpclattice.GetTargetGroupInput{
		TargetGroupIdentifier: aws.String(id),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *VPCLatticeClient) ListServiceNetworks() ([]types.ServiceNetworkSummary, error) {
	serviceNetworks := []types.ServiceNetworkSummary{}
	pageNum := 0

	paginator := vpclattice.NewListServiceNetworksPaginator(c.Client, &vpclattice.ListServiceNetworksInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		serviceNetworks = append(serviceNetworks, out.Items...)
		pageNum++
	}

	return serviceNetworks, nil
}

func (c *VPCLatticeClient) ListServices() ([]types.ServiceSummary, error) {
	services := []types.ServiceSummary{}
	pageNum := 0

	paginator := vpclattice.NewListServicesPaginator(c.Client, &vpclattice.ListServicesInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		services = append(services, out.Items...)
		pageNum++
	}

	return services, nil
}

func (c *VPCLatticeClient) ListTargetGroups() ([]types.TargetGroupSummary, error) {
	targetGroups := []types.TargetGroupSummary{}
	pageNum := 0

	paginator := vpclattice.NewListTargetGroupsPaginator(c.Client, &vpclattice.ListTargetGroupsInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		targetGroups = append(targetGroups, out.Items...)
		pageNum++
	}

	return targetGroups, nil
}
