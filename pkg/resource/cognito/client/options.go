package client

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
	AppClientID  string
	UserPoolID   string
	UserPoolName string

	// Create, Delete
	ClientName string

	// Create
	OAuthScopes  []string
	CallbackUrls []string
}

func NewOptions() (options *Options, createFlags, deleteFlags, getFlags cmd.Flags) {
	options = &Options{
		CommonOptions: resource.CommonOptions{
			Name:                "cognito-domain",
			ClusterFlagDisabled: true,
		},
		CallbackUrls: []string{"http://localhost"},
		OAuthScopes:  []string{"openid"},
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
					if options.UserPoolName == "" {
						return nil
					}

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
	deleteFlags = commonFlags
	getFlags = append(commonFlags,
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "id",
				Description: "get by id instead of name",
				Validate: func(_ *cobra.Command, args []string) error {
					if options.AppClientID != "" && len(args) > 0 {
						return &cmd.ArgumentAndFlagCantBeUsedTogetherError{Arg: "NAME", Flag: "--id"}
					}
					return nil
				},
			},
			Option: &options.AppClientID,
		},
	)

	return
}

func (o *Options) SetName(name string) {
	o.ClientName = name
}
