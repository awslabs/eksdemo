package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type EKSClient struct {
	*eks.Client
}

func NewEKSClient() *EKSClient {
	return &EKSClient{eks.NewFromConfig(GetConfig())}
}

func (c *EKSClient) DescribeAddon(clusterName, addonName string) (*types.Addon, error) {
	result, err := c.Client.DescribeAddon(context.Background(), &eks.DescribeAddonInput{
		AddonName:   aws.String(addonName),
		ClusterName: aws.String(clusterName),
	})

	if err != nil {
		return nil, err
	}

	return result.Addon, nil
}

func (c *EKSClient) DescribeAddonVersions(addonName, version string) ([]types.AddonInfo, error) {
	addons := []types.AddonInfo{}
	pageNum := 0

	input := eks.DescribeAddonVersionsInput{
		KubernetesVersion: aws.String(version),
	}

	if addonName != "" {
		input.AddonName = aws.String(addonName)
	}

	paginator := eks.NewDescribeAddonVersionsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		addons = append(addons, out.Addons...)
		pageNum++
	}

	return addons, nil
}

func (c *EKSClient) DescribeCluster(clusterName string) (*types.Cluster, error) {
	result, err := c.Client.DescribeCluster(context.Background(), &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})

	if err != nil {
		return nil, err
	}

	return result.Cluster, nil
}

func (c *EKSClient) DescribeFargateProfile(clusterName, profileName string) (*types.FargateProfile, error) {
	result, err := c.Client.DescribeFargateProfile(context.Background(), &eks.DescribeFargateProfileInput{
		ClusterName:        aws.String(clusterName),
		FargateProfileName: aws.String(profileName),
	})

	if err != nil {
		return nil, err
	}

	return result.FargateProfile, nil
}

func (c *EKSClient) DescribeNodegroup(clusterName, nodegroupName string) (*types.Nodegroup, error) {
	result, err := c.Client.DescribeNodegroup(context.Background(), &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodegroupName),
	})

	if err != nil {
		return nil, err
	}

	return result.Nodegroup, nil
}

func (c *EKSClient) ListAddons(clusterName string) ([]string, error) {
	addons := []string{}
	pageNum := 0

	paginator := eks.NewListAddonsPaginator(c.Client, &eks.ListAddonsInput{
		ClusterName: aws.String(clusterName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		addons = append(addons, out.Addons...)
		pageNum++
	}

	return addons, nil
}

func (c *EKSClient) ListClusters() ([]string, error) {
	clusters := []string{}
	pageNum := 0

	paginator := eks.NewListClustersPaginator(c.Client, &eks.ListClustersInput{})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, out.Clusters...)
		pageNum++
	}

	return clusters, nil
}

func (c *EKSClient) ListFargateProfiles(clusterName string) ([]string, error) {
	profileNames := []string{}
	pageNum := 0

	paginator := eks.NewListFargateProfilesPaginator(c.Client, &eks.ListFargateProfilesInput{
		ClusterName: aws.String(clusterName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		profileNames = append(profileNames, out.FargateProfileNames...)
		pageNum++
	}

	return profileNames, nil
}

func (c *EKSClient) ListNodegroups(clusterName string) ([]string, error) {
	nodegroupNames := []string{}
	pageNum := 0

	paginator := eks.NewListNodegroupsPaginator(c.Client, &eks.ListNodegroupsInput{
		ClusterName: aws.String(clusterName),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		nodegroupNames = append(nodegroupNames, out.Nodegroups...)
		pageNum++
	}

	return nodegroupNames, nil
}

func (c *EKSClient) UpdateNodegroupConfig(clusterName, nodegroupName string, desired, min, max int) error {
	_, err := c.Client.UpdateNodegroupConfig(context.Background(), &eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(clusterName),
		NodegroupName: aws.String(nodegroupName),
		ScalingConfig: &types.NodegroupScalingConfig{
			DesiredSize: aws.Int32(int32(desired)),
			MinSize:     aws.Int32(int32(min)),
			MaxSize:     aws.Int32(int32(max)),
		},
	})

	return err
}
