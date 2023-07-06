package argo

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

type GuestbookOptions struct {
	resource.CommonOptions
	ArgoCDNamespace string
}

func NewGuestbook() *resource.Resource {
	options := &GuestbookOptions{
		ArgoCDNamespace: "argocd",
		CommonOptions: resource.CommonOptions{
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "guestbook",
			Description: "Argo Guestbook Application",
		},

		CreateFlags: cmd.Flags{
			&cmd.StringFlag{
				CommandFlag: cmd.CommandFlag{
					Name:        "argocd-namespace",
					Description: "namespace where Argo CD is deployed",
				},
				Option: &options.ArgoCDNamespace,
			},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: workspaceYamlTemplate,
			},
		},

		Options: options,
	}
}

const workspaceYamlTemplate = `---
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: guestbook
  namespace: {{ .ArgoCDNamespace }}
  finalizers:
    - resources-finalizer.argocd.argoproj.io
spec:
  destination:
    namespace: {{ .Namespace }}
    server: https://kubernetes.default.svc
  project: default
  source:
    path: guestbook
    repoURL: https://github.com/argoproj/argocd-example-apps
    targetRevision: HEAD
  syncPolicy:
    automated: {}
    syncOptions:
      - CreateNamespace=true
`
