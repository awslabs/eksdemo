package nodegroup

import (
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

// /aws/service/eks/optimized-ami/<eks-version>/amazon-linux-2/recommended/image_id
const eksOptmizedAmiPath = "/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/image_id"

type NodegroupOptions struct {
	*resource.CommonOptions

	AMI             string
	InstanceType    string
	Containerd      bool
	DesiredCapacity int
	MinSize         int
	MaxSize         int
	NodegroupName   string
	OperatingSystem string
	Spot            bool
	SpotvCPUs       int
	SpotMemory      int

	UpdateDesired int
	UpdateMin     int
	UpdateMax     int
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
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "containerd",
				Description: "use containerd runtime",
				Validate: func(cmd *cobra.Command, args []string) error {
					if !options.Containerd {
						return nil
					}
					if v := options.Common().KubernetesVersion; v == "1.23" || v == "1.22" || v == "1.21" {
						return nil
					}
					return fmt.Errorf("%q flag not supported for EKS versions 1.24 and later", "containerd")
				},
			},
			Option: &options.Containerd,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "instance",
				Description: "instance type",
				Shorthand:   "i",
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
					if options.Containerd && options.OperatingSystem != "AmazonLinux2" {
						return fmt.Errorf("%q flag can only be used with %q Operating System", "containerd", "AmazonLinux2")
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
	if !o.Containerd {
		return nil
	}

	param, err := aws.NewSSMClient().GetParameter(fmt.Sprintf(eksOptmizedAmiPath, o.KubernetesVersion))
	if err != nil {
		return fmt.Errorf("ssm failed to lookup EKS Optimized AMI: %w", err)
	}

	o.AMI = awssdk.ToString(param.Value)

	return nil
}

func (o *NodegroupOptions) SetName(name string) {
	o.NodegroupName = name
}
