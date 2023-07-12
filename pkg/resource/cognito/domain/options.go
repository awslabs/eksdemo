package domain

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/userpool"
	"github.com/spf13/cobra"
)

type Options struct {
	resource.CommonOptions
	DomainName string

	// Create, Delete
	UserPoolID   string
	UserPoolName string
}

func NewOptions() (options *Options, createFlags, deleteFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			Name:                "cognito-domain",
			ClusterFlagDisabled: true,
			GetArgumentRequired: true,
		},
	}

	commonFlags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "user-pool-id",
				Description: "id of the user pool",
				Shorthand:   "I",
				Validate: func(cmd *cobra.Command, args []string) error {
					if options.UserPoolID == "" && options.UserPoolName == "" {
						return fmt.Errorf("must include either %q flag or %q flag", "--user-pool-id", "--user-pool-name")
					}
					return nil
				},
			},
			Option: &options.UserPoolID,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "user-pool-name",
				Description: "name of the user pool",
				Shorthand:   "U",
				Validate: func(cmd *cobra.Command, args []string) error {
					up, err := userpool.NewGetter(aws.NewCognitoUserPoolClient()).GetUserPoolByName(options.UserPoolName)
					if err != nil {
						return err
					}
					options.UserPoolID = awssdk.ToString(up.Id)
					return nil
				},
			},
			Option: &options.UserPoolName,
		},
	}

	createFlags = commonFlags

	deleteFlags = cmd.Flags{}

	return
}

func (o *Options) SetName(name string) {
	o.DomainName = name
}
