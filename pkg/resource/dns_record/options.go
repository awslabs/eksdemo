package dns_record

import (
	"fmt"
	"strings"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type DnsRecordOptions struct {
	resource.CommonOptions
	// Shared
	ZoneName string
	// Create
	Type  string
	Value string
	// Delete
	AllRecords bool
	AllTypes   bool
	// Get
	Filter []string
}

func newOptions() (options *DnsRecordOptions, createFlags, deleteFlags, getFlags cmd.Flags) {
	options = &DnsRecordOptions{
		CommonOptions: resource.CommonOptions{
			DeleteArgumentOptional: true,
			ClusterFlagDisabled:    true,
		},
	}

	commonFlags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "zone",
				Description: "hosted zone name",
				Shorthand:   "z",
				Required:    true,
			},
			Option: &options.ZoneName,
		},
	}

	createFlags = append(cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "type",
				Description: "record type",
				Shorthand:   "t",
				Required:    true,
				Validate: func(cmd *cobra.Command, args []string) error {
					if strings.EqualFold(options.Type, "A") {
						options.Type = "A"
						return nil
					}
					if strings.EqualFold(options.Type, "CNAME") {
						options.Type = "CNAME"
						return nil
					}
					if strings.EqualFold(options.Type, "TXT") {
						options.Type = "TXT"
						return nil
					}
					return nil
				},
			},
			Choices: []string{"A", "CNAME", "TXT"},
			Option:  &options.Type,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "value",
				Description: "record value",
				Shorthand:   "v",
				Required:    true,
			},
			Option: &options.Value,
		},
	}, commonFlags...)

	deleteFlags = append(cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all",
				Description: "delete all records types for the record name",
				Shorthand:   "A",
			},
			Option: &options.AllTypes,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "all-records",
				Description: "delete all records in the zone (excluding root records)",
				Validate: func(cmd *cobra.Command, args []string) error {
					if !options.AllRecords && len(args) == 0 {
						return fmt.Errorf("must include either %q argument or %q flag", "NAME", "--all-records")
					}

					if options.AllRecords && len(args) > 0 {
						return fmt.Errorf("%q flag cannot be used with a record name", "--all-records")
					}
					return nil
				},
			},
			Option: &options.AllRecords,
		},
	}, commonFlags...)

	getFlags = append(cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "filter",
				Description: "filter records by types",
			},
			Option: &options.Filter,
		},
	}, commonFlags...)

	return
}
