package instance

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

const amazonLinux2023AmiPath = "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-%s"

type Options struct {
	resource.CommonOptions

	// Create
	AMI                    string
	Count                  int
	InstanceType           string
	KeyName                string
	SubnetID               string
	VolumeSize             int
	supportedArchitectures []types.ArchitectureType
	volumeDeviceName       string

	// Get
	HideTerminated bool
}

func newOptions() (o *Options, createFlags, getFlags cmd.Flags) {
	o = &Options{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled:    true,
			ClusterFlagOptional:    true,
			CreateArgumentOptional: true,
		},
		Count:        1,
		InstanceType: "t2.micro",
	}

	createFlags = cmd.Flags{
		// Instance Type flag MUST come before AMI flag because because AMI flag validation uses
		// o.supportedArchitectures that is set during the validation of the instance-type flag
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "instance-type",
				Description: "instance type",
				Shorthand:   "i",
				Validate: func(_ *cobra.Command, _ []string) error {
					instanceTypes, err := aws.NewEC2Client().DescribeInstanceTypes(
						[]types.Filter{aws.NewEC2InstanceTypeFilter(o.InstanceType)},
					)
					if err != nil {
						return fmt.Errorf("failed to describe instance types: %w", err)
					}

					if len(instanceTypes) != 1 {
						return fmt.Errorf("%q is not a valid instance type in region %q", o.InstanceType, aws.Region())
					}

					o.supportedArchitectures = instanceTypes[0].ProcessorInfo.SupportedArchitectures

					return nil
				},
			},
			Option: &o.InstanceType,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ami",
				Description: "ID of the AMI (defaults to latest AL2023)",
				Validate: func(_ *cobra.Command, _ []string) error {
					if o.AMI != "" {
						return nil
					}

					arch := ""

					for _, sa := range o.supportedArchitectures {
						if sa == types.ArchitectureTypeX8664 || sa == types.ArchitectureTypeArm64 {
							arch = string(sa)
							break
						}
					}

					if arch == "" {
						return fmt.Errorf(
							"%q flag is required for %q", "ami", o.InstanceType,
						)
					}

					// Default to latest AL 2023 AMI
					param, err := aws.NewSSMClient().GetParameter(fmt.Sprintf(amazonLinux2023AmiPath, arch))
					if err != nil {
						return fmt.Errorf(
							"failed to lookup Amazon Linux 2023 AMI for architecture %q: %w", arch, err,
						)
					}
					o.AMI = awssdk.ToString(param.Value)

					return nil
				},
			},
			Option: &o.AMI,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "count",
				Description: "number of instances to launch",
				Shorthand:   "n",
			},
			Option: &o.Count,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "key-name",
				Description: "name of the key pair",
			},
			Option: &o.KeyName,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "subnet",
				Description: "ID of the subnet to launch the instance into",
			},
			Option: &o.SubnetID,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "volume-size",
				Description: "volume size in GiB",
				Shorthand:   "v",
				Validate: func(_ *cobra.Command, _ []string) error {
					if o.VolumeSize == 0 {
						return nil
					}

					amis, err := aws.NewEC2Client().DescribeImages(nil, []string{o.AMI}, nil)
					if err != nil {
						return err
					}

					if len(amis) != 1 {
						return fmt.Errorf("failed to describe AMI %q looking up block device mappings", o.AMI)
					}

					if len(amis[0].BlockDeviceMappings) == 0 {
						return fmt.Errorf("no block device mapping found for AMI %q", o.AMI)
					}

					o.volumeDeviceName = awssdk.ToString(amis[0].BlockDeviceMappings[0].DeviceName)

					return nil
				},
			},
			Option: &o.VolumeSize,
		},
	}

	getFlags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "hide-terminated",
				Description: "show only pending, running, shutting-down, stopping, stopped instances",
				Validate: func(_ *cobra.Command, args []string) error {
					if o.HideTerminated && len(args) > 0 {
						return fmt.Errorf("%q flag cannot be used with ID argument", "hide-terminated")
					}
					return nil
				},
			},
			Option: &o.HideTerminated,
		},
	}
	return
}
