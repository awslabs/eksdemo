package installer

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/helm"
	"github.com/awslabs/eksdemo/pkg/kustomize"
	"github.com/awslabs/eksdemo/pkg/template"

	"helm.sh/helm/v3/pkg/chart"
)

type ManifestInstaller struct {
	AppName           string
	DryRun            bool
	KustomizeTemplate template.Template
	ResourceTemplate  template.Template
}

func (i *ManifestInstaller) Install(options application.Options) error {
	resources, err := i.ResourceTemplate.Render(options)
	if err != nil {
		return err
	}

	var yaml string

	if i.KustomizeTemplate != nil {
		kustomization, err := i.KustomizeTemplate.Render(options)
		if err != nil {
			return err
		}

		yaml, err = kustomize.Kustomize(resources, kustomization)
		if err != nil {
			return err
		}
		yaml = "---\n" + yaml
	} else {
		yaml = resources
	}

	if i.DryRun {
		fmt.Println("\nManifest Installer Dry Run:")
		fmt.Println(yaml)
		return nil
	}

	chart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    i.AppName,
			Version: "n/a",
			Type:    "application",
		},
		Templates: []*chart.File{
			{
				Name: "main",
				Data: []byte(yaml),
			},
		},
	}

	o := options.Common()
	if o.Version == "" {
		o.Version = "n/a"
	}

	h := helm.Helm{
		AppVersion:  o.Version,
		Namespace:   o.Namespace,
		ReleaseName: i.AppName,
		ValuesFile:  "",
	}

	return h.Install(chart, o.KubeContext())
}

func (i *ManifestInstaller) SetDryRun() {
	i.DryRun = true
}

func (i *ManifestInstaller) Type() application.InstallerType {
	return application.ManifestInstaller
}

func (i *ManifestInstaller) Uninstall(options application.Options) error {
	o := options.Common()

	fmt.Printf("Checking status of application: %s, in namespace: %s\n", i.AppName, o.Namespace)
	if _, err := helm.Status(o.KubeContext(), i.AppName, o.Namespace); err != nil {
		return err
	}

	fmt.Println("Status validated. Uninstalling...")
	return helm.Uninstall(o.KubeContext(), i.AppName, o.Namespace)
}
