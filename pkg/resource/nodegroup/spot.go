package nodegroup

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/eksctl"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewSpotResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "nodegroup-spot",
			Description: "Managed Nodegroup with Spot Instances",
			Aliases:     []string{"spot", "ngspot", "ng-spot"},
			Args:        []string{"NAME"},
		},

		Manager: &eksctl.ResourceManager{
			Resource: "nodegroup",
			ConfigTemplate: &template.TextTemplate{
				Template: eksctl.EksctlHeader + EksctlTemplate,
			},
			ApproveDelete: true,
		},
	}

	res.Options, res.CreateFlags = NewSpotOptions()

	return res
}

func NewSpotOptions() (options *NodegroupOptions, flags cmd.Flags) {
	options, flags, _ = NewOptions()
	options.Spot = true

	spotFlags := cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "vcpus",
				Description: "use instance types with specified vCPUs",
			},
			Option: &options.SpotvCPUs,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "mem",
				Description: "use instance types with specified memory",
			},
			Option: &options.SpotMemory,
		},
	}

	flags = append(flags.Remove("instance"), spotFlags...)
	return
}
