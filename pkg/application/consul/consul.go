package consul

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

// https://developer.hashicorp.com/consul/docs/k8s/installation/install
// https://developer.hashicorp.com/consul/tutorials/kubernetes/kubernetes-eks-aws
// https://artifacthub.io/packages/helm/hashicorp/consul

func NewApp() *application.Application {
	app := &application.Application{
		Command: cmd.Command{
			Name:        "consul",
			Description: "HashiCorp Consul Service-Mesh",
		},

		Installer: &installer.HelmInstaller{
			ChartName:     "consul",
			ReleaseName:   "consul",
			RepositoryURL: "https://helm.releases.hashicorp.com",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate,
			},
		},
	}

	app.Options, app.Flags = newOptions()

	return app

}

const valuesTemplate = `---
global:
  # The main enabled/disabled setting.
  # If true, servers, clients, Consul DNS and the Consul UI will be enabled.
  enabled: true
  # The prefix used for all resources created in the Helm chart.
  name: null
  # The name of the datacenter that the agents should register as.
  datacenter: {{ .Datacenter }}
  # Enables TLS across the cluster to verify authenticity of the Consul servers and clients.
  tls:
    enabled: true
  # Enables ACLs across the cluster to secure access to data and APIs.
  acls:
  # If true, automatically manage ACL tokens and policies for all Consul components.
    manageSystemACLs: true

# Configures values that configure the Consul server cluster.
server:
  enabled: true
  # The number of server agents to run. This determines the fault tolerance of the cluster.
  replicas: {{ .Replicas }}

  # Contains values that configure the Consul UI.
{{ if .EnableUI }}
ui:
  enabled: true
  # Registers a Kubernetes Service for the Consul UI as a LoadBalancer.
  service:
    type: LoadBalancer
    annotations: |
      service.beta.kubernetes.io/aws-load-balancer-type: "external"
      service.beta.kubernetes.io/aws-load-balancer-scheme: "internet-facing"
      service.beta.kubernetes.io/aws-load-balancer-nlb-target-type: "ip"
      service.beta.kubernetes.io/aws-load-balancer-attributes: "load_balancing.cross_zone.enabled=true"
{{ end }}

# Configures and installs the automatic Consul Connect sidecar injector.
{{ if .EnableMesh }}
connectInject:
  enabled: true
    # Enables Consul on Kubernetes to manage the CRDs used for Gateway API.
    # Setting this to true will install the CRDs used for the Gateway API when Consul on Kubernetes is installed.
    # These CRDs can clash with existing Gateway API CRDs if they are already installed in your cluster.
    # If this setting is false, you will need to install the Gateway API CRDs manually.
    {{ if .EnableAPIGW }}
  apiGateway:
    manageExternalCRDs: true
  managedGatewayClass:
    serviceType: LoadBalancer
    {{ end }}
{{ end }}
`
