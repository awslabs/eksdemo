package secretsstore

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// Docs:    https://secrets-store-csi-driver.sigs.k8s.io/
// GitHub:  https://github.com/kubernetes-sigs/secrets-store-csi-driver
// Helm:    https://github.com/kubernetes-sigs/secrets-store-csi-driver/tree/main/charts/secrets-store-csi-driver
// Repo:    registry.k8s.io/csi-secrets-store/driver
// Version: Latest is v1.4.4 (as of 8/6/24)

func NewApp() *application.Application {
	options, flags := newOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "secrets",
			Name:        "store-csi-driver",
			Description: "Integrates secrets stores with K8s via a CSI volume",
			Aliases:     []string{"store-csi", "csi-driver", "csi"},
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "secrets-store-csi-driver",
			ReleaseName:   "secrets-store-csi-driver",
			RepositoryURL: "https://kubernetes-sigs.github.io/secrets-store-csi-driver/charts",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},

		Options: options,
	}
}

// https://github.com/kubernetes-sigs/secrets-store-csi-driver/blob/main/charts/secrets-store-csi-driver/values.yaml
const valuesTemplate = `---
linux:
  image:
    tag: {{ .Version }}
  crds:
    image:
      tag: {{ .Version }}
syncSecret:
  enabled: true
{{- if .RotateEnabled }}
enableSecretRotation: {{ .RotateEnabled }}
rotationPollInterval: {{ .RotationPollInterval }}
{{- end }}
`
