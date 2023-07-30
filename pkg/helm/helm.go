package helm

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/awslabs/eksdemo/pkg/printer"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/postrender"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/strvals"
	"sigs.k8s.io/yaml"
)

var debug bool

type Helm struct {
	AppVersion    string
	ChartName     string
	ChartVersion  string
	Namespace     string
	PostRenderer  postrender.PostRenderer
	ReleaseName   string
	RepositoryURL string
	Wait          bool

	SetValues  []string
	ValuesFile string
}

func Init(helmDebug bool) {
	debug = helmDebug
}

func initialize(kubeContext, namespace string) (*action.Configuration, error) {
	// Hack to work around https://github.com/helm/helm/issues/7430
	_ = os.Setenv("HELM_KUBECONTEXT", kubeContext)
	_ = os.Setenv("HELM_NAMESPACE", namespace)
	settings := cli.New()

	// Initialize the action configuration
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, "secret", log.Printf); err != nil {
		return nil, fmt.Errorf("failed to initialize helm action config: %w", err)
	}
	return actionConfig, nil
}

func (h *Helm) DownloadChart() (*chart.Chart, error) {
	getters := getter.All(&cli.EnvSettings{})

	u, err := url.Parse(h.RepositoryURL)
	if err != nil {
		return nil, err
	}

	var chartPath string
	if u.Scheme == "oci" {
		chartPath = h.RepositoryURL + ":" + h.ChartVersion
	} else {
		// Find Chart
		chartPath, err = repo.FindChartInRepoURL(h.RepositoryURL, h.ChartName, h.ChartVersion, "", "", "", getters)
		if err != nil {
			return nil, err
		}
	}

	fmt.Printf("Downloading Chart: %s\n", chartPath)
	g, err := getters.ByScheme(u.Scheme)
	if err != nil {
		return nil, err
	}

	// Download chart archive into memory
	data, err := g.Get(chartPath)

	// If ECR Public is returning a 403 Forbidden error, then log out and try again
	// https://docs.aws.amazon.com/AmazonECR/latest/public/public-troubleshooting.html
	if h.isECRPublicAuthError(err) {
		fmt.Println("ECR Public is returning a 403 Forbidden error. Logging out and trying again...")

		registryClient, e := registry.NewClient()
		if e != nil {
			return nil, fmt.Errorf("failed to create registry client: %w", e)
		}

		e = registryClient.Logout("public.ecr.aws")
		if e != nil {
			return nil, fmt.Errorf("failed to log out of ECR Public: %w", e)
		}

		ociGetter, e := getter.NewOCIGetter()
		if e != nil {
			return nil, e
		}
		data, err = ociGetter.Get(chartPath)
	}
	if err != nil {
		return nil, err
	}

	// Decompress the archive
	files, err := loader.LoadArchiveFiles(data)
	if err != nil {
		return nil, err
	}

	// Load the chart
	chart, err := loader.LoadFiles(files)
	if err != nil {
		return nil, err
	}
	return chart, nil
}

func (h *Helm) Install(chart *chart.Chart, kubeContext string) error {
	// Parse the values file
	values := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(h.ValuesFile), &values); err != nil {
		return fmt.Errorf("failed to parse values file: %w", err)
	}

	for _, v := range h.SetValues {
		if err := strvals.ParseInto(v, values); err != nil {
			return fmt.Errorf("failed parsing --set data: %w", err)
		}
	}

	if debug {
		fmt.Println("\nHelm Debug, Final Values:")
		fmt.Println("---")
		err := printer.EncodeYAML(os.Stdout, values)
		if err != nil {
			return fmt.Errorf("failed to print final values: %w", err)
		}
		fmt.Println()
	}

	actionConfig, err := initialize(kubeContext, h.Namespace)
	if err != nil {
		return err
	}

	// Configure the install options
	instAction := action.NewInstall(actionConfig)
	instAction.Namespace = h.Namespace
	instAction.ReleaseName = h.ReleaseName
	instAction.CreateNamespace = true
	instAction.IsUpgrade = true
	instAction.PostRenderer = h.PostRenderer
	instAction.Wait = h.Wait
	instAction.Timeout = 300 * time.Second
	chart.Metadata.AppVersion = h.AppVersion

	// Install the chart
	fmt.Println("Helm installing...")
	rel, err := instAction.Run(chart, values)
	if err != nil {
		return fmt.Errorf("helm install failed: %s", err)
	}

	fmt.Printf("Using chart version %q, installed %q version %q in namespace %q\n",
		rel.Chart.Metadata.Version, rel.Name, rel.Chart.Metadata.AppVersion, rel.Namespace)

	// Print the Chart NOTES
	if len(rel.Info.Notes) > 0 {
		fmt.Printf("NOTES:\n%s\n", strings.TrimSpace(rel.Info.Notes))
	}

	return nil
}

func (h *Helm) isECRPublicAuthError(err error) bool {
	return err != nil &&
		strings.HasPrefix(h.RepositoryURL, "oci://public.ecr.aws") &&
		strings.HasSuffix(err.Error(), "403 Forbidden")
}

func List(kubeContext string) ([]*release.Release, error) {
	actionConfig, err := initialize(kubeContext, "")
	if err != nil {
		return nil, err
	}

	client := action.NewList(actionConfig)
	client.AllNamespaces = true

	releases, err := client.Run()
	if (err) != nil {
		return nil, err
	}

	return releases, nil
}

func Status(kubeContext, releaseName, namespace string) (string, error) {
	actionConfig, err := initialize(kubeContext, namespace)
	if err != nil {
		return "", err
	}
	status := action.NewStatus(actionConfig)

	rel, err := status.Run(releaseName)
	if (err) != nil {
		return "", err
	}

	// strip chart metadata from the output
	rel.Chart = nil

	return "", nil
}

func Uninstall(kubeContext, releaseName, namespace string) error {
	actionConfig, err := initialize(kubeContext, namespace)
	if err != nil {
		return err
	}
	uninstall := action.NewUninstall(actionConfig)

	// Uninstall the chart
	_, err = uninstall.Run(releaseName)
	if err != nil {
		return fmt.Errorf("failed uninstalling chart: %w", err)
	}

	return nil
}
