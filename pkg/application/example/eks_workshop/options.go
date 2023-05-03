package eks_workshop

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type EksWorkshopOptions struct {
	application.ApplicationOptions

	CrystalReplicas  int
	FrontendReplicas int
	NodeReplicas     int
}

func NewOptions() (options *EksWorkshopOptions, flags cmd.Flags) {
	options = &EksWorkshopOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:                    "eks-workshop",
			DisableServiceAccountFlag:    true,
			DisableVersionFlag:           true,
			ExposeIngressAndLoadBalancer: true,
		},
		CrystalReplicas:  3,
		FrontendReplicas: 3,
		NodeReplicas:     3,
	}

	flags = cmd.Flags{
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "crystal-replicas",
				Description: "number of replicas for the Crystal deployment",
			},
			Option: &options.CrystalReplicas,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "frontend-replicas",
				Description: "number of replicas for the Frontend deployment",
			},
			Option: &options.FrontendReplicas,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "nodejs-replicas",
				Description: "number of replicas for the Node.js deployment",
			},
			Option: &options.NodeReplicas,
		},
	}
	return
}
