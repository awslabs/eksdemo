package velero

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/s3_bucket"
)

type VeleroOptions struct {
	application.ApplicationOptions

	PluginVersion string
	*s3_bucket.BucketOptions
}

func newOptions() (options *VeleroOptions, flags cmd.Flags) {
	options = &VeleroOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "velero",
			ServiceAccount: "velero-server",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "5.0.2",
				Latest:        "v1.11.1",
				PreviousChart: "5.0.2",
				Previous:      "v1.11.1",
			},
		},
		PluginVersion: "v1.7.1",
		BucketOptions: &s3_bucket.BucketOptions{
			CommonOptions: resource.CommonOptions{
				Name: "velero-s3-bucket",
			},
		},
	}
	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "plugin-version",
				Description: "aws plugin version",
			},
			Option: &options.PluginVersion,
		},
	}
	return
}

func (o *VeleroOptions) PreDependencies(application.Action) error {
	o.BucketOptions.BucketName = fmt.Sprintf("eksdemo-%s-velero", o.Account)
	return nil
}
