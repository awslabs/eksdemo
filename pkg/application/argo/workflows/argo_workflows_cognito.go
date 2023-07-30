package workflows

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/goutils"
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/client"
	"github.com/awslabs/eksdemo/pkg/resource/cognito/userpool"
	"github.com/awslabs/eksdemo/pkg/template"
)

const appClientName = "argo-workflows"

type CognitoOptions struct {
	*Options
	UserPoolOptions *userpool.Options

	ClientID     string
	ClientSecret string
	OAuthScopes  []string
	UserPoolID   string
}

func NewAppWithCognito() *application.Application {
	options, flags := newCognitoOptions()

	return &application.Application{
		Command: cmd.Command{
			Parent:      "argo",
			Name:        "workflows-cognito",
			Description: "Workflow engine for Kubernetes using Cognito for authentication",
		},

		Dependencies: []*resource.Resource{
			userpool.NewWithOptions(options.UserPoolOptions),
		},

		Flags: flags,

		Installer: &installer.HelmInstaller{
			ChartName:     "argo-workflows",
			ReleaseName:   "argo-workflows-cognito",
			RepositoryURL: "https://argoproj.github.io/argo-helm",
			ValuesTemplate: &template.TextTemplate{
				Template: valuesTemplate + cognitoValuesTemplate,
			},
			PostRenderKustomize: &template.TextTemplate{
				Template: postRenderKustomize,
			},
		},

		Options: options,
	}
}

func newCognitoOptions() (options *CognitoOptions, flags cmd.Flags) {
	workflowOptions, flags := newOptions()

	workflowOptions.AuthMode = "sso"
	workflowOptions.IngressOnly = true

	options = &CognitoOptions{
		Options:         workflowOptions,
		OAuthScopes:     []string{"email"},
		UserPoolOptions: &userpool.Options{},
	}
	return
}

// https://github.com/argoproj/argo-helm/blob/main/charts/argo-workflows/values.yaml
const cognitoValuesTemplate = `
  sso:
    enabled: true
    issuer: https://cognito-idp.{{ .Region }}.amazonaws.com/{{ .UserPoolID }}
    clientId:
      name: argo-server-sso
      key: client-id
    clientSecret:
      name: argo-server-sso
      key: client-secret
    redirectUrl: https://{{ .IngressHost }}/oauth2/callback
    rbac:
      enabled: false
    scopes: {{ .OAuthScopes }}
extraObjects:
  - apiVersion: v1
    kind: Secret
    metadata:
      name: argo-server-sso
    data:
      client-id: {{ .ClientID | b64enc }}
      client-secret: {{ .ClientSecret | b64enc}}
    type: Opaque
`

// Workaround for https://github.com/argoproj/argo-helm/issues/2159
const postRenderKustomize = `---
resources:
- manifest.yaml
patches:
# Add service account permission to the argo server clusterrole
- patch: |-
    - op: add
      path: /rules/-
      value:
        apiGroups:
        - ""
        resources:
        - serviceaccounts
        verbs:
        - get
        - list
        - watch
  target:
    group: rbac.authorization.k8s.io
    kind: ClusterRole
    name: argo-server
    version: v1
`

func (o *CognitoOptions) userPoolName() string {
	return fmt.Sprintf("%s-%s", o.ClusterName, "argo")
}

func (o *CognitoOptions) PreDependencies(application.Action) error {
	o.UserPoolOptions.UserPoolName = o.userPoolName()
	return nil
}

func (o *CognitoOptions) PreInstall() error {
	if o.DryRun {
		fmt.Println("\nPreInstall Dry Run:")
		fmt.Println("1. Check if the user pool has a domain configured.")
		fmt.Println("2. If no domain is configured, create a Cognito domain.")
		fmt.Println("3. Create a new Cognito app client.")
		return nil
	}

	cognitoClient := aws.NewCognitoUserPoolClient()

	userPool, err := userpool.NewGetter(cognitoClient).GetUserPoolByName(o.UserPoolOptions.UserPoolName)
	if err != nil {
		return err
	}

	o.UserPoolID = awssdk.ToString(userPool.Id)

	if userPool.Domain == nil {
		fmt.Printf("User Pool %q does not have a domain configured.\n", o.UserPoolOptions.UserPoolName)

		rand, err := goutils.CryptoRandomAlphaNumeric(8)
		if err != nil {
			return fmt.Errorf("failed to generate random string: %w", err)
		}

		prefix := fmt.Sprintf("eksdemo-argo-%s", strings.ToLower(rand))
		fmt.Printf("Creating new cognito domain with prefix %q...", prefix)

		_, err = cognitoClient.CreateUserPoolDomain(prefix, o.UserPoolID)
		if err != nil {
			return err
		}
		fmt.Println("done")
	}

	_, err = client.NewGetter(cognitoClient).GetAppClientByName(appClientName, o.UserPoolID)
	if err == nil {
		return fmt.Errorf("app client %q already exists for user pool %q, please delete it and try again",
			appClientName, o.UserPoolID)
	}

	// Return the error if it's anything except "not found" error
	var notFoundErr *resource.NotFoundByNameError
	if err != nil && !errors.As(err, &notFoundErr) {
		return err
	}

	fmt.Printf("Creating new app client %q...", appClientName)
	appClient, err := cognitoClient.CreateUserPoolClient(
		append(o.OAuthScopes, "openid"),
		[]string{fmt.Sprintf("https://%s/oauth2/callback", o.IngressHost)},
		appClientName,
		o.UserPoolID,
	)
	if err != nil {
		return err
	}
	fmt.Println("done")

	o.ClientID = awssdk.ToString(appClient.ClientId)
	o.ClientSecret = awssdk.ToString(appClient.ClientSecret)

	return nil
}
