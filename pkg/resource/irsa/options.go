package irsa

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/template"
	"github.com/spf13/cobra"
)

type IrsaOptions struct {
	resource.CommonOptions

	PolicyType
	Policy []string

	// Used for flags
	PolicyARNs        []string
	PolicyDocTemplate template.Template
}

type PolicyType int

const (
	None PolicyType = iota
	PolicyARNs
	PolicyDocument
)

func addOptions(res *resource.Resource) *resource.Resource {
	options := &IrsaOptions{
		CommonOptions: resource.CommonOptions{
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	res.Options = options

	res.CreateFlags = cmd.Flags{
		&cmd.StringSliceFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "attach-arns",
				Description: "ARNs",
				Validate: func(cmd *cobra.Command, args []string) error {
					if len(options.PolicyARNs) == 0 {
						return nil
					}

					if len(options.Policy) > 0 {
						return fmt.Errorf("can only use one policy flag")
					}

					options.PolicyType = PolicyARNs
					options.Policy = options.PolicyARNs

					return nil
				},
			},
			Option: &options.PolicyARNs,
		},
	}

	return res
}

func (o *IrsaOptions) ClusterOIDCProvider() (string, error) {
	issuer := aws.ToString(o.Cluster.Identity.Oidc.Issuer)

	slices := strings.Split(issuer, "//")
	if len(slices) < 2 {
		return "", fmt.Errorf("failed to parse Cluster OIDC Provider URL")
	}

	return slices[1], nil
}

func (o *IrsaOptions) IrsaAnnotation() string {
	return fmt.Sprintf("eks.amazonaws.com/role-arn: arn:%s:iam::%s:role/%s", o.Partition, o.Account, o.RoleName())
}

func (o *IrsaOptions) IsPolicyDocument(t PolicyType) bool {
	return t == PolicyDocument
}

func (o *IrsaOptions) IsPolicyARN(t PolicyType) bool {
	return t == PolicyARNs
}

func (o *IrsaOptions) RoleName() string {
	return o.TruncateUnique(64, fmt.Sprintf("eksdemo.%s.%s.%s", o.ClusterName, o.Namespace, o.ServiceAccount))
}

func (o *IrsaOptions) SetName(name string) {
	o.ServiceAccount = name
}
