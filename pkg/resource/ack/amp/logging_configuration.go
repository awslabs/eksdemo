package amp

import (
	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/manifest"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/log_group"
	"github.com/awslabs/eksdemo/pkg/template"
)

type LoggingConfigurationOptions struct {
	resource.CommonOptions

	Alias        string
	LogGroupName string
	LogGroupArn  string
}

func NewLoggingConfigurationResource() *resource.Resource {
	res := &resource.Resource{
		Command: cmd.Command{
			Name:        "amp-logging-configuration",
			Description: "AMP Logging Configuration",
			Aliases:     []string{"amp-logging-config", "amp-logging", "amp-log"},
			Args:        []string{"NAME"},
		},

		Manager: &manifest.ResourceManager{
			Template: &template.TextTemplate{
				Template: loggingYamlTemplate,
			},
		},
	}

	options := &LoggingConfigurationOptions{
		CommonOptions: resource.CommonOptions{
			Name:          "ack-amp-logging-configuration",
			Namespace:     "default",
			NamespaceFlag: true,
		},
	}

	flags := cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "alias",
				Description: "workspace alias to configure logging",
				Required:    true,
			},
			Option: &options.Alias,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "log-group",
				Description: "name of CloudWatch Log Group",
				Required:    true,
			},
			Option: &options.LogGroupName,
		},
	}

	res.Options = options
	res.CreateFlags = flags

	return res
}

func (o *LoggingConfigurationOptions) PreCreate() error {
	lg, err := log_group.NewGetter(aws.NewCloudwatchlogsClient()).GetLogGroupByName(o.LogGroupName)
	if err != nil {
		return err
	}

	o.LogGroupArn = awssdk.ToString(lg.Arn)

	return nil
}

const loggingYamlTemplate = `---
apiVersion: prometheusservice.services.k8s.aws/v1alpha1
kind: LoggingConfiguration
metadata:
  namespace: {{ .Namespace }}
  name: {{ .Name }}
spec:
  logGroupARN: {{ .LogGroupArn }}
  workspaceRef:
    from:
      name: {{ .Alias }}
`
