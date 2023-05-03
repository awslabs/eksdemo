package eks

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/subnet"
	"github.com/awslabs/eksdemo/pkg/template"
)

type FargateProfileOptions struct {
	resource.CommonOptions
	FargateNamespace string
	Subnets          []string
}

func NewFargateProfileResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "eks-fargate-profile",
			Description: "Fargate Profile",
			Aliases:     []string{"eks-fargate", "fargate-profile", "fargate", "fp"},
			Args:        []string{"NAME"},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: subnetYamlTemplate,
			},
		},
	}

	options := &FargateProfileOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-eks-fargate-profile",
			Namespace:     "default",
			NamespaceFlag: true,
		},
		FargateNamespace: "default",
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "fargate-namespace",
				Description: "namespace selector to run pods on fargate",
			},
			Option: &options.FargateNamespace,
		},
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "subnets",
				Description: "subnets for fargate pods (defaults to all private subnets)",
			},
			Option: &options.Subnets,
		},
	}

	res.Options = options
	res.CreateFlags = flags

	return res
}

func (o *FargateProfileOptions) PreCreate() error {
	if len(o.Subnets) > 0 {
		return nil
	}

	subnets, err := subnet.NewGetter(aws.NewEC2Client()).GetPrivateSubnetsForCluster(o.Cluster)
	if err != nil {
		return err
	}

	if len(subnets) == 0 {
		return fmt.Errorf("subnet autodiscovery failed, use --subnets flag")
	}

	for _, s := range subnets {
		o.Subnets = append(o.Subnets, awssdk.ToString(s.SubnetId))
	}

	return nil
}

const subnetYamlTemplate = `---
apiVersion: eks.services.k8s.aws/v1alpha1
kind: FargateProfile
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  clusterName: {{ .ClusterName }}
  name: {{ .Name }}
  podExecutionRoleARN: arn:aws:iam::{{ .Account }}:role/eksdemo.{{ .ClusterName }}.fargate-pod-execution-role
  selectors:
  - namespace: {{ .FargateNamespace }}
  subnets:
{{- range .Subnets }}
  - {{ . }}
{{- end }}
`
