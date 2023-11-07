package coredumphandler

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/s3_bucket"
)

type Options struct {
	application.ApplicationOptions

	*s3_bucket.BucketOptions
	IncludeCrioExe bool
}

func newOptions() (options *Options, flags cmd.Flags) {
	options = &Options{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "core-dump-handler",
			ServiceAccount: "core-dump-handler",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "v8.10.0",
				Latest:        "v8.10.0",
				PreviousChart: "v8.10.0",
				Previous:      "v8.10.0",
			},
		},
		BucketOptions: &s3_bucket.BucketOptions{
			CommonOptions: resource.CommonOptions{
				Name: "coredumphandler-s3-bucket",
			},
		},
	}
	flags = cmd.Flags{
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "dont-include-crio-exe",
				Description: "don't include CRI-O executable",
			},
			Option: &options.IncludeCrioExe,
		},
	}
	return
}

func (o *Options) PreDependencies(application.Action) error {
	o.BucketOptions.BucketName = fmt.Sprintf("eksdemo-%s-coredumphandler", o.Account)
	return nil
}
