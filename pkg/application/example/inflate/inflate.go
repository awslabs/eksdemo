package inflate

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Manifest: https://karpenter.sh/v0.27.5/getting-started/getting-started-with-karpenter/#scale-up-deployment
// Repo:     https://gallery.ecr.aws/eks-distro/kubernetes/pause

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "example",
			Name:        "inflate",
			Description: "Example Deployment for Karpenter autoscaling",
		},

		Flags: flags,

		Installer: &installer.ManifestInstaller{
			AppName: "example-inflate",
			ResourceTemplate: &template.TextTemplate{
				Template: manifestTemplate,
			},
		},

		Options: options,
	}
}

// https://karpenter.sh/docs/getting-started/getting-started-with-karpenter/#6-scale-up-deployment
const manifestTemplate = `---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inflate
  namespace: {{ .Namespace }}
spec:
  replicas: {{ .Replicas }}
  selector:
    matchLabels:
      app: inflate
  template:
    metadata:
      labels:
        app: inflate
    spec:
      terminationGracePeriodSeconds: 0
      securityContext:
        runAsUser: 1000
        runAsGroup: 3000
        fsGroup: 2000
      containers:
      - name: inflate
        image: public.ecr.aws/eks-distro/kubernetes/pause:3.7
        resources:
          requests:
            cpu: 1
        securityContext:
          allowPrivilegeEscalation: false
{{- if .OnDemand }}
      nodeSelector:
        karpenter.sh/capacity-type: on-demand
{{- end }}
{{- if .Spread }}
      topologySpreadConstraints:
      - maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
        labelSelector:
          matchLabels:
            app: inflate
{{- end }}
`
