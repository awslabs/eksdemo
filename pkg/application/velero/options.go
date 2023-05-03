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
				LatestChart:   "2.30.1",
				Latest:        "v1.9.0",
				PreviousChart: "2.30.1",
				Previous:      "v1.9.0",
			},
		},
		PluginVersion: "v1.5.0",
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
