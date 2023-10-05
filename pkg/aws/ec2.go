package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Client struct {
	*ec2.Client
}

func NewEC2Client() *EC2Client {
	return &EC2Client{ec2.NewFromConfig(GetConfig())}
}

func NewEC2DescriptionFilter(description string) types.Filter {
	return types.Filter{
		Name:   aws.String("description"),
		Values: []string{description},
	}
}

func NewEC2ElasticIpFilter(eipId string) types.Filter {
	return types.Filter{
		Name:   aws.String("allocation-id"),
		Values: []string{eipId},
	}
}

func NewEC2InstanceIdFilter(instanceId string) types.Filter {
	return types.Filter{
		Name:   aws.String("instance-id"),
		Values: []string{instanceId},
	}
}

func NewEC2InstanceStateFilter(states []string) types.Filter {
	return types.Filter{
		Name:   aws.String("instance-state-name"),
		Values: states,
	}
}

func NewEC2InstanceTypeFilter(instanceType string) types.Filter {
	return types.Filter{
		Name:   aws.String("instance-type"),
		Values: []string{instanceType},
	}
}

func NewEC2InternetGatewayFilter(internetGatewayId string) types.Filter {
	return types.Filter{
		Name:   aws.String("internet-gateway-id"),
		Values: []string{internetGatewayId},
	}
}

func NewEC2InternetGatewayVpcFilter(vpcId string) types.Filter {
	return types.Filter{
		Name:   aws.String("attachment.vpc-id"),
		Values: []string{vpcId},
	}
}

func NewEC2NatGatewayFilter(natGatewayId string) types.Filter {
	return types.Filter{
		Name:   aws.String("nat-gateway-id"),
		Values: []string{natGatewayId},
	}
}

func NewEC2NetworkAclFilter(networkAclId string) types.Filter {
	return types.Filter{
		Name:   aws.String("network-acl-id"),
		Values: []string{networkAclId},
	}
}

func NewEC2NetworkInterfaceFilter(networkInterfaceId string) types.Filter {
	return types.Filter{
		Name:   aws.String("network-interface-id"),
		Values: []string{networkInterfaceId},
	}
}

func NewEC2NetworkInterfaceInstanceFilter(instanceId string) types.Filter {
	return types.Filter{
		Name:   aws.String("attachment.instance-id"),
		Values: []string{instanceId},
	}
}

func NewEC2NetworkInterfaceIpFilter(ipAddress string) types.Filter {
	return types.Filter{
		Name:   aws.String("addresses.private-ip-address"),
		Values: []string{ipAddress},
	}
}

func NewEC2PrefixListFilter(prefixListId string) types.Filter {
	return types.Filter{
		Name:   aws.String("prefix-list-id"),
		Values: []string{prefixListId},
	}
}

func NewEC2RouteTableFilter(routeTableId string) types.Filter {
	return types.Filter{
		Name:   aws.String("route-table-id"),
		Values: []string{routeTableId},
	}
}

func NewEC2SecurityGroupFilter(securityGroupId string) types.Filter {
	return types.Filter{
		Name:   aws.String("group-id"),
		Values: []string{securityGroupId},
	}
}

func NewEC2SecurityGroupRuleFilter(securityGroupRuleId string) types.Filter {
	return types.Filter{
		Name:   aws.String("security-group-rule-id"),
		Values: []string{securityGroupRuleId},
	}
}

func NewEC2SubnetFilter(subnetId string) types.Filter {
	return types.Filter{
		Name:   aws.String("subnet-id"),
		Values: []string{subnetId},
	}
}

func NewEC2TagKeyFilter(tagKey string) types.Filter {
	return types.Filter{
		Name:   aws.String("tag-key"),
		Values: []string{tagKey},
	}
}

func NewEC2VpcEndpointFilter(vpcEndpointId string) types.Filter {
	return types.Filter{
		Name:   aws.String("vpc-endpoint-id"),
		Values: []string{vpcEndpointId},
	}
}

func NewEC2VpcFilter(vpcId string) types.Filter {
	return types.Filter{
		Name:   aws.String("vpc-id"),
		Values: []string{vpcId},
	}
}

// Adds or overwrites only the specified tags for the specified Amazon EC2 resource or resources.
// When you specify an existing tag key, the value is overwritten with the new value.
// Each resource can have a maximum of 50 tags. Each tag consists of a key and optional value.
// Tag keys must be unique per resource.
func (c *EC2Client) CreateTags(resources []string, tags map[string]string) error {
	_, err := c.Client.CreateTags(context.Background(), &ec2.CreateTagsInput{
		Resources: resources,
		Tags:      toEC2Tags(tags),
	})

	return err
}

// Deletes a security group.
func (c *EC2Client) DeleteSecurityGroup(id string) error {
	_, err := c.Client.DeleteSecurityGroup(context.Background(), &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(id),
	})

	return err
}

// Deletes the specified EBS volume. The volume must be in the available state (not attached to an instance).
func (c *EC2Client) DeleteVolume(id string) error {
	_, err := c.Client.DeleteVolume(context.Background(), &ec2.DeleteVolumeInput{
		VolumeId: aws.String(id),
	})

	return err
}

// Describes the specified Elastic IP addresses or all of your Elastic IP addresses.
func (c *EC2Client) DescribeAddresses(filters []types.Filter) ([]types.Address, error) {
	input := ec2.DescribeAddressesInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	result, err := c.Client.DescribeAddresses(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return result.Addresses, nil
}

// Describes the Availability Zones, Local Zones, and Wavelength Zones that are available to you.
// If there is an event impacting a zone, you can use this request to view the state and any provided messages for that zone.
func (c *EC2Client) DescribeAvailabilityZones(name string, all bool) ([]types.AvailabilityZone, error) {
	filters := []types.Filter{}
	input := ec2.DescribeAvailabilityZonesInput{
		AllAvailabilityZones: aws.Bool(all),
	}

	if name != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("zone-name"),
			Values: []string{name},
		})
	}

	if len(filters) > 0 {
		input.Filters = filters
	}

	result, err := c.Client.DescribeAvailabilityZones(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return result.AvailabilityZones, nil
}

// Describes the specified instances or all instances.
func (c *EC2Client) DescribeInstances(filters []types.Filter) ([]types.Reservation, error) {
	reservations := []types.Reservation{}
	pageNum := 0

	input := ec2.DescribeInstancesInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeInstancesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, out.Reservations...)
		pageNum++
	}

	return reservations, nil
}

// Describes the details of the instance types that are offered in a location.
// The results can be filtered by the attributes of the instance types.
func (c *EC2Client) DescribeInstanceTypes(filters []types.Filter) ([]types.InstanceTypeInfo, error) {
	instanceTypes := []types.InstanceTypeInfo{}
	pageNum := 0

	input := ec2.DescribeInstanceTypesInput{
		MaxResults: aws.Int32(100),
	}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeInstanceTypesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		instanceTypes = append(instanceTypes, out.InstanceTypes...)
		pageNum++
	}

	return instanceTypes, nil
}

// Describes one or more of your internet gateways.
func (c *EC2Client) DescribeInternetGateways(filters []types.Filter) ([]types.InternetGateway, error) {
	internetGateways := []types.InternetGateway{}
	pageNum := 0

	input := ec2.DescribeInternetGatewaysInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeInternetGatewaysPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		internetGateways = append(internetGateways, out.InternetGateways...)
		pageNum++
	}

	return internetGateways, nil
}

// Describes your managed prefix lists and any AWS-managed prefix lists.
// To view the entries for your prefix list, use GetManagedPrefixListEntries.
func (c *EC2Client) DescribeManagedPrefixLists(filters []types.Filter) ([]types.ManagedPrefixList, error) {
	prefixLists := []types.ManagedPrefixList{}
	pageNum := 0

	input := ec2.DescribeManagedPrefixListsInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeManagedPrefixListsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		prefixLists = append(prefixLists, out.PrefixLists...)
		pageNum++
	}

	return prefixLists, nil
}

// Describes one or more of your NAT gateways.
func (c *EC2Client) DescribeNATGateways(filters []types.Filter) ([]types.NatGateway, error) {
	nats := []types.NatGateway{}
	pageNum := 0

	input := ec2.DescribeNatGatewaysInput{}
	if len(filters) > 0 {
		input.Filter = filters
	}

	paginator := ec2.NewDescribeNatGatewaysPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		nats = append(nats, out.NatGateways...)
		pageNum++
	}

	return nats, nil
}

// Describes one or more of your network ACLs.
func (c *EC2Client) DescribeNetworkAcls(filters []types.Filter) ([]types.NetworkAcl, error) {
	nacls := []types.NetworkAcl{}
	pageNum := 0

	input := ec2.DescribeNetworkAclsInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeNetworkAclsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		nacls = append(nacls, out.NetworkAcls...)
		pageNum++
	}

	return nacls, nil
}

// Describes one or more of your network interfaces.
func (c *EC2Client) DescribeNetworkInterfaces(filters []types.Filter) ([]types.NetworkInterface, error) {
	networkInterfaces := []types.NetworkInterface{}
	pageNum := 0

	input := ec2.DescribeNetworkInterfacesInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeNetworkInterfacesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		networkInterfaces = append(networkInterfaces, out.NetworkInterfaces...)
		pageNum++
	}

	return networkInterfaces, nil
}

// Describes one or more of your route tables.
func (c *EC2Client) DescribeRouteTables(filters []types.Filter) ([]types.RouteTable, error) {
	routeTables := []types.RouteTable{}
	pageNum := 0

	input := ec2.DescribeRouteTablesInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeRouteTablesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		routeTables = append(routeTables, out.RouteTables...)
		pageNum++
	}

	return routeTables, nil
}

// Describes one or more of your security group rules.
func (c *EC2Client) DescribeSecurityGroupRules(id, securityGroupId string) ([]types.SecurityGroupRule, error) {
	filters := []types.Filter{}
	securityGroupRules := []types.SecurityGroupRule{}
	pageNum := 0

	if id != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("security-group-rule-id"),
			Values: []string{id},
		})
	}

	if securityGroupId != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("group-id"),
			Values: []string{securityGroupId},
		})
	}

	input := ec2.DescribeSecurityGroupRulesInput{}

	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeSecurityGroupRulesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		securityGroupRules = append(securityGroupRules, out.SecurityGroupRules...)
		pageNum++
	}

	return securityGroupRules, nil
}

// Describes the specified security groups or all of your security groups.
func (c *EC2Client) DescribeSecurityGroups(filters []types.Filter, groupIds []string) ([]types.SecurityGroup, error) {
	securityGroups := []types.SecurityGroup{}
	pageNum := 0

	input := ec2.DescribeSecurityGroupsInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	if len(groupIds) > 0 {
		input.GroupIds = groupIds
	}

	paginator := ec2.NewDescribeSecurityGroupsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		securityGroups = append(securityGroups, out.SecurityGroups...)
		pageNum++
	}

	return securityGroups, nil
}

// Describes one or more of your subnets.
func (c *EC2Client) DescribeSubnets(filters []types.Filter) ([]types.Subnet, error) {
	subnets := []types.Subnet{}
	pageNum := 0

	input := ec2.DescribeSubnetsInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeSubnetsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		subnets = append(subnets, out.Subnets...)
		pageNum++
	}

	return subnets, nil
}

// Describes the specified tags for your EC2 resources.
func (c *EC2Client) DescribeTags(resources, tagsFilter []string) ([]types.TagDescription, error) {
	tags := []types.TagDescription{}
	pageNum := 0

	filters := []types.Filter{
		{
			Name:   aws.String("resource-id"),
			Values: resources,
		},
	}

	if len(tagsFilter) > 0 {
		filters = append(filters, types.Filter{
			Name:   aws.String("key"),
			Values: tagsFilter,
		})
	}

	paginator := ec2.NewDescribeTagsPaginator(c.Client, &ec2.DescribeTagsInput{
		Filters: filters,
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		tags = append(tags, out.Tags...)
		pageNum++
	}

	return tags, nil
}

// Describes the specified EBS volumes or all of your EBS volumes.
func (c *EC2Client) DescribeVolumes(id string) ([]types.Volume, error) {
	filters := []types.Filter{}
	volumes := []types.Volume{}
	pageNum := 0

	if id != "" {
		filters = append(filters, types.Filter{
			Name:   aws.String("volume-id"),
			Values: []string{id},
		})
	}

	input := ec2.DescribeVolumesInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeVolumesPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		volumes = append(volumes, out.Volumes...)
		pageNum++
	}

	return volumes, nil
}

// Describes your VPC endpoints.
func (c *EC2Client) DescribeVpcEndpoints(filters []types.Filter) ([]types.VpcEndpoint, error) {
	vpcEndpoints := []types.VpcEndpoint{}
	pageNum := 0

	input := ec2.DescribeVpcEndpointsInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeVpcEndpointsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		vpcEndpoints = append(vpcEndpoints, out.VpcEndpoints...)
		pageNum++
	}

	return vpcEndpoints, nil
}

// Describes one or more of your VPCs.
func (c *EC2Client) DescribeVpcs(filters []types.Filter) ([]types.Vpc, error) {
	vpcs := []types.Vpc{}
	pageNum := 0

	input := ec2.DescribeVpcsInput{}
	if len(filters) > 0 {
		input.Filters = filters
	}

	paginator := ec2.NewDescribeVpcsPaginator(c.Client, &input)

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		vpcs = append(vpcs, out.Vpcs...)
		pageNum++
	}

	return vpcs, nil
}

// Gets information about the entries for a specified managed prefix list.
func (c *EC2Client) GetManagedPrefixListEntries(prefixListId string) ([]types.PrefixListEntry, error) {
	pageNum := 0
	entries := []types.PrefixListEntry{}

	paginator := ec2.NewGetManagedPrefixListEntriesPaginator(c.Client, &ec2.GetManagedPrefixListEntriesInput{
		PrefixListId: aws.String(prefixListId),
	})

	for paginator.HasMorePages() && pageNum < maxPages {
		out, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		entries = append(entries, out.Entries...)
		pageNum++
	}

	return entries, nil
}

// Shuts down the specified instances.
// This operation is idempotent; if you terminate an instance more than once, each call succeeds.
func (c *EC2Client) TerminateInstances(id string) error {
	_, err := c.Client.TerminateInstances(context.Background(), &ec2.TerminateInstancesInput{
		InstanceIds: []string{id},
	})

	return err
}

func toEC2Tags(tags map[string]string) (ec2tags []types.Tag) {
	for k, v := range tags {
		ec2tags = append(ec2tags, types.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}
	return
}
