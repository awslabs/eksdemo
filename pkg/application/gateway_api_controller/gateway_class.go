package gateway_api_controller

import (
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
)

func gatewayClass() *resource.Resource {
	res := &resource.Resource{
		Options: &resource.CommonOptions{
			Name: "amazon-vpc-lattice-gateway-class",
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: yamlTemplate,
			},
		},
	}
	return res
}

// https://github.com/aws/aws-application-networking-k8s/blob/main/examples/gatewayclass.yaml
const yamlTemplate = `---
# Create a new Gateway Class for AWS VPC lattice provider
apiVersion: gateway.networking.k8s.io/v1beta1
kind: GatewayClass
metadata:
  name: amazon-vpc-lattice
spec:
  controllerName: application-networking.k8s.aws/gateway-api-controller
`
