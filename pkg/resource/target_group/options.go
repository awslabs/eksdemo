package target_group

import (
	"fmt"
	"strings"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type TargeGroupOptions struct {
	resource.CommonOptions

	LoadBalancerName string
	Protocol         string
	TargetType       string
	VpcId            string
}

func newOptions() (options *TargeGroupOptions, createFlags, getFlags cmd.Flags) {
	options = &TargeGroupOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
		Protocol:   "http",
		TargetType: "instance",
	}

	createFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "protocol",
				Description: "protocol for routing traffic to the targets",
				Shorthand:   "p",
				Validate: func(cmd *cobra.Command, args []string) error {
					options.Protocol = strings.ToUpper(options.Protocol)
					return nil
				},
			},
			Choices: []string{"http", "https", "tcp", "tls", "udp", "tcp_udp", "geneve"},
			Option:  &options.Protocol,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "target-type",
				Description: "target type",
				Shorthand:   "t",
				Validate: func(cmd *cobra.Command, args []string) error {
					options.TargetType = strings.ToLower(options.TargetType)
					return nil
				},
			},
			Choices: []string{"ip", "instance", "lambda", "alb"},
			Option:  &options.TargetType,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "vpc-id",
				Description: "VPC to create the target group",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.ClusterName == "" && options.VpcId == "" {
						return fmt.Errorf("must include either %q or %q flag", "--cluster", "--vpc-id")
					}

					if options.ClusterName != "" && options.VpcId != "" {
						return fmt.Errorf("%q flag and %q flag can not be used together", "--cluster", "--vpc-id")
					}

					return nil
				},
			},
			Option: &options.VpcId,
		},
	}

	getFlags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "load-balancer",
				Description: "filter by Load Balancer name",
				Shorthand:   "L",
			},
			Option: &options.LoadBalancerName,
		},
	}

	return
}
