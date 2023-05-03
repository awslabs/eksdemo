package keycloak_amg

import (
	"context"
	"fmt"
	"net/http"
	"time"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/amg_workspace"
	"github.com/go-resty/resty/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const AmgAliasSuffix = `keycloak-amg`
const samlMetadataPath = `realms/eksdemo/protocol/saml/descriptor`

type KeycloakOptions struct {
	application.ApplicationOptions

	AdminPassword   string
	AmgWorkspaceUrl string
	amgWorkspaceId  string
	*amg_workspace.AmgOptions
}

func NewOptions() (options *KeycloakOptions, flags cmd.Flags) {
	options = &KeycloakOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "9.6.8",
				Latest:        "18.0.2",
				PreviousChart: "9.2.10",
				Previous:      "18.0.0",
			},
			ExposeIngressOnly: true,
			Namespace:         "keycloak",
			ServiceAccount:    "keycloak",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "admin-pass",
				Description: "Keycloak admin password",
				Required:    true,
				Shorthand:   "P",
			},
			Option: &options.AdminPassword,
		},
	}
	return
}

func (o *KeycloakOptions) PreDependencies(application.Action) error {
	o.AmgOptions.WorkspaceName = fmt.Sprintf("%s-%s", o.ClusterName, AmgAliasSuffix)
	return nil
}

func (o *KeycloakOptions) PreInstall() error {
	grafanaGetter := amg_workspace.NewGetter(aws.NewGrafanaClient())

	workspace, err := grafanaGetter.GetAmgByName(o.AmgOptions.WorkspaceName)
	if err != nil {
		if o.DryRun {
			o.amgWorkspaceId = "<AMG Workspace ID>"
			return nil
		}
		return fmt.Errorf("failed to lookup AMG URL to use in Helm chart: %w", err)
	}

	o.amgWorkspaceId = awssdk.ToString(workspace.Id)
	o.AmgWorkspaceUrl = awssdk.ToString(workspace.Endpoint)

	return nil
}

func (o *KeycloakOptions) PostInstall(_ string, _ []*resource.Resource) error {
	if o.DryRun {
		fmt.Println("Postinstall will update AMG Workspace to complete SAML configuration")
		return nil
	}

	fmt.Print("Waiting for Keycloak SAML metadata URL to become active...")

	var metadataUrl string
	if o.IngressHost == "" {
		k8sclient, err := kubernetes.Client(o.KubeContext())
		if err != nil {
			return err
		}

		svc, err := k8sclient.CoreV1().Services(o.Namespace).Get(context.Background(), keycloakReleaseName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if len(svc.Status.LoadBalancer.Ingress) == 0 {
			return fmt.Errorf("failed to get Service Load Balancer address")
		}

		metadataUrl = fmt.Sprintf("http://%s/%s", svc.Status.LoadBalancer.Ingress[0].Hostname, samlMetadataPath)
	} else {
		metadataUrl = fmt.Sprintf("https://%s/%s", o.IngressHost, samlMetadataPath)
	}

	_, err := resty.New().
		SetLogger(NewLogger()).
		SetRetryCount(10).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				if err != nil {
					// retry on "no such host" error when using TLS as we wait for DNS
					return true
				}
				// retry as we wait for the ALB healthcheck for the target group
				return r.StatusCode() == http.StatusServiceUnavailable
			},
		).
		R().Get(metadataUrl)

	if err != nil {
		fmt.Println()
		return fmt.Errorf("%w\n\nTo finish configuration, update AMG with the SAML metadata URL: %s", err, metadataUrl)
	}
	fmt.Println("done")
	fmt.Printf("Updating AMG with Keyclock SAML Metadata URL to complete SAML configuration\n")

	err = aws.NewGrafanaClient().UpdateWorkspaceAuthentication(o.amgWorkspaceId, metadataUrl)
	if err != nil {
		fmt.Println("Metadata URL is: " + metadataUrl)
		return err
	}
	fmt.Printf("Amazon Managed Grafana available at: https://%s\n", o.AmgWorkspaceUrl)

	return nil
}
