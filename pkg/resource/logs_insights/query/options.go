package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type QueryOptions struct {
	resource.CommonOptions

	Fields      []string
	FilterField string
	FilterLike  string
	FilterRaw   string
	Limit       int
	LogGroup    string
	Query       string
	SortField   string
	SortOrder   string
	TimeRange   time.Duration
}

func NewOptions() (options *QueryOptions, createFlags cmd.Flags) {
	options = &QueryOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
		Fields:      []string{"@timestamp", "@message", "@logStream"},
		FilterField: "@message",
		Limit:       10000,
		SortField:   "@timestamp",
		SortOrder:   "asc",
		TimeRange:   24 * time.Hour,
	}

	createFlags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "fields",
				Description: "fields to include in query results",
			},
			Option: &options.Fields,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "filter-field",
				Description: "log field to filter, when using --filter-like flag",
				Shorthand:   "F",
				Validate: func(cmd *cobra.Command, args []string) error {
					if cmd.Flags().Changed("filter-field") && options.FilterLike == "" {
						return fmt.Errorf("%q flag requires %q flag", "--filter-field", "--filter-like")
					}
					return nil
				},
			},
			Option: &options.FilterField,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "filter-like",
				Description: "substring or /regular expression/ for \"like\" keyword phrase",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.FilterLike == "" {
						return nil
					}

					// If FilterLike is a substring (not a regular expression), put quotes around it
					if !strings.HasPrefix(options.FilterLike, "/") && !strings.HasSuffix(options.FilterLike, "/") {
						options.FilterLike = fmt.Sprintf("%q", options.FilterLike)
					}
					return nil
				},
			},
			Option: &options.FilterLike,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "filter-raw",
				Description: "raw filter command, must be quoted properly",
			},
			Option: &options.FilterRaw,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "limit",
				Description: "max number of log events to return",
				Shorthand:   "L",
			},
			Option: &options.Limit,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "log-group",
				Description: "log group to query",
				Shorthand:   "l",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.ClusterName == "" && options.LogGroup == "" {
						return fmt.Errorf("must include either %q or %q flag", "--cluster", "--log-group")
					}

					if options.ClusterName != "" && options.LogGroup != "" {
						return fmt.Errorf("%q flag and %q flag can not be used together", "--cluster", "--log-group")
					}

					return nil
				},
			},
			Option: &options.LogGroup,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "sort",
				Description: "display in asc or desc order",
			},
			Choices: []string{"asc", "desc"},
			Option:  &options.SortOrder,
		},
		&cmd.DurationFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "time",
				Description: "time range to query logs",
				Shorthand:   "t",
			},
			Option: &options.TimeRange,
		},
	}
	return
}
