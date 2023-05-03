package log_group

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

func newOptions() (options *resource.CommonOptions, getFlags, deleteFlags cmd.Flags) {
	options = &resource.CommonOptions{
		DeleteArgumentOptional: true,
		ClusterFlagDisabled:    true,
	}

	getClusterFlag := options.NewClusterFlag(resource.Get, false)
	getClusterFlag.Validate = func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && options.ClusterName != "" {
			return fmt.Errorf("%q argument and %q flag can not be used together", "NAME", "--cluster")
		}

		return nil
	}

	deleteClusterFlag := options.NewClusterFlag(resource.Delete, false)
	deleteClusterFlag.Validate = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && options.ClusterName == "" {
			return fmt.Errorf("must include either %q argument or %q flag", "NAME", "--cluster")
		}

		if options.Name == "" {
			options.Name = LogGroupNameForClusterName(options.ClusterName)
		}

		return nil
	}

	getFlags = cmd.Flags{getClusterFlag}
	deleteFlags = cmd.Flags{deleteClusterFlag}

	return
}
