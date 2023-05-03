package amp_workspace

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, flags := NewOptions()
	res := NewResourceWithOptions(options)
	res.DeleteFlags = flags

	return res
}

func NewResourceWithOptions(options *AmpWorkspaceOptions) *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amp-workspace",
			Description: "Amazon Managed Prometheus Workspace",
			Aliases:     []string{"amp"},
			Args:        []string{"ALIAS"},
		},

		Getter: &Getter{},

		Manager: &Manager{},
	}

	res.Options = options

	return res
}
