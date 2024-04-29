package installer

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/helm"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/kustomize"
	"github.com/awslabs/eksdemo/pkg/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type HelmInstaller struct {
	ChartName           string
	DryRun              bool
	PostRenderKustomize template.Template
	PVCLabels           map[string]string
	ReleaseName         string
	RepositoryURL       string
	ValuesTemplate      template.Template
	Wait                bool

	options application.Options
}

func (i *HelmInstaller) Install(options application.Options) error {
	valuesFile, err := i.ValuesTemplate.Render(options)
	if err != nil {
		return err
	}

	if i.DryRun {
		fmt.Println("\nHelm Installer Dry Run:")
		NewHelmPrinter(i, options, valuesFile).PrintTable(os.Stdout)

		if i.PostRenderKustomize != nil {
			kustomization, err := i.PostRenderKustomize.Render(options)
			if err != nil {
				return err
			}
			fmt.Println("Helm Installer Post Render Kustomize Dry Run:")
			fmt.Println(kustomization)
		}
		return nil
	}

	helm := &helm.Helm{
		AppVersion:    options.Common().Version,
		ChartName:     i.ChartName,
		ChartVersion:  options.Common().ChartVersion,
		Namespace:     options.Common().Namespace,
		ReleaseName:   i.ReleaseName,
		RepositoryURL: i.RepositoryURL,
		Wait:          i.Wait,
		SetValues:     options.Common().SetValues,
		ValuesFile:    valuesFile,
	}

	if i.PostRenderKustomize != nil {
		i.options = options
		helm.PostRenderer = i
	}

	chart, err := helm.DownloadChart()
	if err != nil {
		return fmt.Errorf("failed to download chart: %w", err)
	}

	return helm.Install(chart, options.KubeContext())
}

func (i *HelmInstaller) SetDryRun() {
	i.DryRun = true
}

func (i *HelmInstaller) Type() application.InstallerType {
	return application.HelmInstaller
}

func (i *HelmInstaller) Uninstall(options application.Options) error {
	o := options.Common()

	fmt.Printf("Checking status of Helm release: %s, in namespace: %s\n", i.ReleaseName, o.Namespace)
	if _, err := helm.Status(o.KubeContext(), i.ReleaseName, o.Namespace); err != nil {
		return err
	}

	fmt.Println("Status validated. Uninstalling...")
	err := helm.Uninstall(o.KubeContext(), i.ReleaseName, o.Namespace)
	if err != nil {
		return err
	}

	if len(i.PVCLabels) == 0 {
		return nil
	}

	// Delete any leftover PVCs as `helm uninstall` won't delete them
	// https://github.com/helm/helm/issues/5156
	client, err := kubernetes.Client(o.KubeContext())
	if err != nil {
		return fmt.Errorf("failed creating kubernetes client: %w", err)
	}

	selector := labels.NewSelector()

	for k, v := range i.PVCLabels {
		req, err := labels.NewRequirement(k, selection.Equals, []string{v})
		if err != nil {
			return err
		}
		selector = selector.Add(*req)
	}

	fmt.Printf("Deleting PVCs with labels: %s\n", selector.String())

	return client.CoreV1().PersistentVolumeClaims(o.Namespace).DeleteCollection(context.Background(),
		metav1.DeleteOptions{},
		metav1.ListOptions{
			LabelSelector: selector.String(),
		},
	)
}

// PostRender
func (h *HelmInstaller) Run(renderedManifests *bytes.Buffer) (modifiedManifests *bytes.Buffer, err error) {
	kustomization, err := h.PostRenderKustomize.Render(h.options)
	if err != nil {
		return nil, err
	}

	yaml, err := kustomize.Kustomize(renderedManifests.String(), kustomization)
	if err != nil {
		return nil, err
	}

	return bytes.NewBufferString(yaml), nil
}
