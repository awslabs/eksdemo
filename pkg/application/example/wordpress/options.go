package wordpress

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type WordpressOptions struct {
	application.ApplicationOptions

	StorageClass      string
	WordpressPassword string
}

func NewOptions() (options *WordpressOptions, flags cmd.Flags) {
	options = &WordpressOptions{
		ApplicationOptions: application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "15.0.13",
				Latest:        "6.0.1",
				PreviousChart: "15.0.4",
				Previous:      "6.0.0",
			},
			ExposeIngressAndLoadBalancer: true,
			Namespace:                    "wordpress",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "storage-class",
				Description: "StorageClass for WordPress and MariaDB Persistent Volumes",
			},
			Option: &options.StorageClass,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "wordpress-pass",
				Description: "WordPress admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.WordpressPassword,
		},
	}
	return
}
