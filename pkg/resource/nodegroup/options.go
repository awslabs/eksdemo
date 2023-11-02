package nodegroup

import (
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2/recommended/image_id
const eksOptmizedAmiPath = "/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id"

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2-arm64/recommended/image_id
const eksOptmizedArmAmiPath = "/aws/service/eks/optimized-ami/%s/amazon-linux-2-arm64/recommended/image_id"

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2-gpu/recommended/image_id
const eksOptmizedGpuAmiPath = "/aws/service/eks/optimized-ami/%s/amazon-linux-2-gpu/recommended/image_id"

type NodegroupOptions struct {
	*resource.CommonOptions

	AMI              string
	InstanceType     string
	IsClusterPrivate bool
	DesiredCapacity  int
	MinSize          int
	MaxSize          int
	NodegroupName    string
	NoTaints         bool
	OperatingSystem  string
	Spot             bool
	SpotvCPUs        int
	SpotMemory       int
	Taints           []Taint

	UpdateDesired int
	UpdateMin     int
	UpdateMax     int
}

type Taint struct {
	Key    string
	Value  string
	Effect string
}

func NewOptions() (options *NodegroupOptions, createFlags, updateFlags cmd.Flags) {
	options = &NodegroupOptions{
		CommonOptions:   &resource.CommonOptions{},
		InstanceType:    "t3.large",
		DesiredCapacity: 1,
		MinSize:         0,
		MaxSize:         10,
		OperatingSystem: "AmazonLinux2",
		SpotvCPUs:       2,
		SpotMemory:      4,
	}

	createFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "instance",
				Description: "instance type",
				Shorthand:   "i",
				Validate: func(cmd *cobra.Command, args []string) error {
					options.InstanceType = strings.ToLower(options.InstanceType)
					return nil
				},
			},
			Option: &options.InstanceType,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "max",
				Description: "max nodes",
			},
			Option: &options.MaxSize,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "min",
				Description: "min nodes",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.MinSize >= options.MaxSize {
						return fmt.Errorf("min nodes must be less than max nodes")
					}
					return nil
				},
			},
			Option: &options.MinSize,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nodes",
				Description: "desired number of nodes",
				Shorthand:   "N",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.DesiredCapacity > options.MaxSize {
						options.MaxSize = options.DesiredCapacity
					}
					if options.DesiredCapacity < options.MinSize {
						options.MinSize = options.DesiredCapacity
					}
					return nil
				},
			},
			Option: &options.DesiredCapacity,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "no-taints",
				Description: "don't taint nodes with GPUs or Neuron cores",
			},
			Option: &options.NoTaints,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "os",
				Description: "Operating System",
				Validate: func(cmd *cobra.Command, args []string) error {
					if strings.EqualFold(options.OperatingSystem, "AmazonLinux2") {
						options.OperatingSystem = "AmazonLinux2"
					}
					if strings.EqualFold(options.OperatingSystem, "Bottlerocket") {
						options.OperatingSystem = "Bottlerocket"
					}
					if strings.EqualFold(options.OperatingSystem, "Ubuntu2004") {
						options.OperatingSystem = "Ubuntu2004"
					}
					if strings.EqualFold(options.OperatingSystem, "Ubuntu1804") {
						options.OperatingSystem = "Ubuntu1804"
					}
					return nil
				},
			},
			Option:  &options.OperatingSystem,
			Choices: []string{"AmazonLinux2", "Bottlerocket", "Ubuntu2004", "Ubuntu1804"},
		},
	}

	updateFlags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "max",
				Description: "max nodes",
			},
			Option: &options.UpdateMax,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "min",
				Description: "min nodes",
			},
			Option: &options.UpdateMin,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nodes",
				Description: "desired number of nodes",
				Shorthand:   "N",
			},
			Option: &options.UpdateDesired,
		},
	}

	return
}

func (o *NodegroupOptions) PreCreate() error {
	instanceTypes, err := aws.NewEC2Client().DescribeInstanceTypes(
		[]types.Filter{aws.NewEC2InstanceTypeFilter(o.InstanceType)},
	)
	if err != nil {
		return fmt.Errorf("failed to describe instance types: %w", err)
	}

	if len(instanceTypes) != 1 {
		return fmt.Errorf("%q is not a valid instance type in region %q", o.InstanceType, o.Region)
	}

	var isGraviton, isNeuron, isNvidia bool
	instType := strings.Split(o.InstanceType, ".")[0]

	switch {
	case strings.HasPrefix(instType, "g"),
		strings.HasPrefix(instType, "p"):

		isNvidia = true

	case strings.HasPrefix(instType, "inf"),
		strings.HasPrefix(instType, "trn"):

		isNeuron = true

	case strings.HasSuffix(instType, "g"),
		strings.HasSuffix(instType, "gd"),
		strings.HasSuffix(instType, "gn"),
		strings.HasSuffix(instType, "gen"):

		isGraviton = true
	}

	if isNeuron && !o.NoTaints {
		o.Taints = append(o.Taints, Taint{Key: "aws.amazon.com/neuron", Effect: "NoSchedule"})
	}

	if isNvidia && !o.NoTaints {
		o.Taints = append(o.Taints, Taint{Key: "nvidia.com/gpu", Effect: "NoSchedule"})
	}

	// AMI Lookup is currently only for Amazon Linux 2 EKS Optimized AMI and clusters that aren't fully private
	if o.OperatingSystem != "AmazonLinux2" || o.IsClusterPrivate {
		return nil
	}

	ssmClient := aws.NewSSMClient()

	switch {
	case instType == "g5g":
		return fmt.Errorf("%q instance type is not supported with the EKS optimized Amazon Linux AMI", "G5g")

	case isNeuron, isNvidia:
		param, err := ssmClient.GetParameter(fmt.Sprintf(eksOptmizedGpuAmiPath, o.KubernetesVersion))
		if err != nil {
			return fmt.Errorf("failed to lookup EKS optimized accelerated AMI for instance type %s: %w", o.InstanceType, err)
		}

		o.AMI = awssdk.ToString(param.Value)

	case isGraviton:
		param, err := ssmClient.GetParameter(fmt.Sprintf(eksOptmizedArmAmiPath, o.KubernetesVersion))
		if err != nil {
			return fmt.Errorf("failed to lookup EKS optimized ARM AMI for instance type %s: %w", o.InstanceType, err)
		}

		o.AMI = awssdk.ToString(param.Value)

	default:
		param, err := ssmClient.GetParameter(fmt.Sprintf(eksOptmizedAmiPath, o.KubernetesVersion))
		if err != nil {
			return fmt.Errorf("failed to lookup EKS optimized AMI for instance type %s: %w", o.InstanceType, err)
		}

		o.AMI = awssdk.ToString(param.Value)
	}

	return nil
}

func (o *NodegroupOptions) SetName(name string) {
	o.NodegroupName = name
}
