package nginx

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://kubernetes.github.io/ingress-nginx/
// GitHub:  https://github.com/kubernetes/ingress-nginx
// Helm:    https://github.com/kubernetes/ingress-nginx/tree/main/charts/ingress-nginx
// Repo:    registry.k8s.io/ingress-nginx/controller
// Version: Latest is Chart 4.7.0 and App v1.8.0 (as of 5/30/23)

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Parent:      "ingress",
			Name:        "nginx",
			Description: "Ingress NGINX Controller",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "ingress-nginx",
			ReleaseName:   "ingress-nginx",
			RepositoryURL: "https://kubernetes.github.io/ingress-nginx",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}
	app.Options, app.Flags = newOptions()

	return app
}

// Service annoations from: https://kubernetes.github.io/ingress-nginx/deploy/#aws
// Added `service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing` annotation
// Because the new AWS LB Controller webhook will default to an internal load balancer
const valuesTemplate = `---
controller:
  image:
    tag: {{ .Version }}
  replicaCount: {{ .Replicas }}
  service:
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-backend-protocol: tcp
      service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled: "true"
      service.beta.kubernetes.io/aws-load-balancer-scheme: internet-facing
      service.beta.kubernetes.io/aws-load-balancer-type: nlb
    externalTrafficPolicy: Local
serviceAccount:
  name: {{ .ServiceAccount }}
`
