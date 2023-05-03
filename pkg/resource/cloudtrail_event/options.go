package cloudtrail_event

import (
	"fmt"
	"strings"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type CloudtrailEventOptions struct {
	resource.CommonOptions

	Ids          bool
	Insights     bool
	Name         string
	ResourceName string
	ResourceType string
	Source       string
	Username     string
}

func newOptions() (o *CloudtrailEventOptions, flags cmd.Flags) {
	o = &CloudtrailEventOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagDisabled: true,
		},
	}

	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "insights",
				Description: "show only Insights events",
			},
			Option: &o.Insights,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "name",
				Description: "filter by event name",
				Shorthand:   "N",
				Validate: func(cmd *cobra.Command, args []string) error {
					if o.Name != "" &&
						(len(args) > 0 || o.ResourceName != "" || o.ResourceType != "" || o.Source != "" || o.Username != "") {
						return fmt.Errorf("only 1 of the following is supported %q argument or %q, %q, %q, %q, %q flags",
							"id", "name", "resource-name", "resource-type", "source", "username")
					}
					return nil
				},
			},
			Option: &o.Name,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "resource-name",
				Description: "filter by resource name",
				Shorthand:   "R",
				Validate: func(cmd *cobra.Command, args []string) error {
					if o.ResourceName != "" &&
						(len(args) > 0 || o.Name != "" || o.ResourceType != "" || o.Source != "" || o.Username != "") {
						return fmt.Errorf("only 1 of the following is supported %q argument or %q, %q, %q, %q, %q flags",
							"id", "name", "resource-name", "resource-type", "source", "username")
					}
					return nil
				},
			},
			Option: &o.ResourceName,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "resource-type",
				Description: "filter by resource type",
				Shorthand:   "T",
				Validate: func(cmd *cobra.Command, args []string) error {
					if o.ResourceType != "" &&
						(len(args) > 0 || o.Name != "" || o.ResourceName != "" || o.Source != "" || o.Username != "") {
						return fmt.Errorf("only 1 of the following is supported %q argument or %q, %q, %q, %q, %q flags",
							"id", "name", "resource-name", "resource-type", "source", "username")
					}
					return nil
				},
			},
			Option: &o.ResourceType,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "source",
				Description: "filter by event source (can leave off \".amazonaws.com\")",
				Shorthand:   "S",
				Validate: func(cmd *cobra.Command, args []string) error {
					if o.Source != "" &&
						(len(args) > 0 || o.Name != "" || o.ResourceName != "" || o.ResourceType != "" || o.Username != "") {
						return fmt.Errorf("only 1 of the following is supported %q argument or %q, %q, %q, %q, %q flags",
							"id", "name", "resource-name", "resource-type", "source", "username")
					}

					if o.Source != "" && !strings.HasSuffix(o.Source, ".amazonaws.com") {
						o.Source += ".amazonaws.com"
					}
					return nil
				},
			},
			Option: &o.Source,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "username",
				Description: "filter by username",
				Shorthand:   "U",
				Validate: func(cmd *cobra.Command, args []string) error {
					if o.Username != "" &&
						(len(args) > 0 || o.Name != "" || o.ResourceName != "" || o.ResourceType != "" || o.Source != "") {
						return fmt.Errorf("only 1 of the following is supported %q argument or %q, %q, %q, %q, %q flags",
							"id", "name", "resource-name", "resource-type", "source", "username")
					}
					return nil
				},
			},
			Option: &o.Username,
		},
	}
	return
}
