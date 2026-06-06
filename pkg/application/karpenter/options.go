package karpenter

import (
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource/ssm/parameter"
	"github.com/spf13/cobra"
)

type KarpenterOptions struct {
	application.ApplicationOptions

	AMIFamily        string
	AMISelectorIDs   []string
	ConsolidateAfter string
	EnableSpotToSpot bool
	ExpireAfter      string
	Replicas         int
}

func newOptions() (options *KarpenterOptions, flags cmd.Flags) {
	options = &KarpenterOptions{
		ApplicationOptions: application.ApplicationOptions{
			Namespace:      "karpenter",
			ServiceAccount: "karpenter",
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "1.5.0",
				Latest:        "1.5.0",
				PreviousChart: "1.0.5",
				Previous:      "1.0.5",
			},
		},
		AMIFamily:        "AL2",
		ConsolidateAfter: "1m",
		ExpireAfter:      "720h",
		Replicas:         1,
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "ami-family",
				Description: "node class AMI family",
				Shorthand:   "a",
				Validate: func(_ *cobra.Command, _ []string) error {
					eksVersion := awssdk.ToString(options.Common().Cluster.Version)
					ssm := parameter.NewGetter(aws.NewSSMClient())

					if strings.EqualFold(options.AMIFamily, "AL2") {
						options.AMIFamily = "AL2"

						al2AMI, err := ssm.GetEKSOptimizedAL2AMI(eksVersion)
						if err != nil {
							return err
						}
						al2Arm64AMI, err := ssm.GetEKSOptimizedAL2Arm64AMI(eksVersion)
						if err != nil {
							return err
						}

						options.AMISelectorIDs = append(options.AMISelectorIDs, al2AMI, al2Arm64AMI)

						return nil
					}
					if strings.EqualFold(options.AMIFamily, "AL2023") {
						options.AMIFamily = "AL2023"

						al2023AMI, err := ssm.GetEKSOptimizedAL2023AMI(eksVersion)
						if err != nil {
							return err
						}
						al2023Arm64AMI, err := ssm.GetEKSOptimizedAL2023Arm64AMI(eksVersion)
						if err != nil {
							return err
						}

						options.AMISelectorIDs = append(options.AMISelectorIDs, al2023AMI, al2023Arm64AMI)
						return nil
					}
					if strings.EqualFold(options.AMIFamily, "Bottlerocket") {
						options.AMIFamily = "Bottlerocket"

						bottlerocketAMI, err := ssm.GetBottlerocketAMI(eksVersion)
						if err != nil {
							return err
						}

						bottlerocketArm64AMI, err := ssm.GetBottlerocketArm64AMI(eksVersion)
						if err != nil {
							return err
						}

						options.AMISelectorIDs = append(options.AMISelectorIDs, bottlerocketAMI, bottlerocketArm64AMI)
					}
					return nil
				},
			},
			Option:  &options.AMIFamily,
			Choices: []string{"AL2", "AL2023", "Bottlerocket"},
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "consolidate-after",
				Description: "time after a pod is scheduled/removed before considering the node consolidatable",
			},
			Option: &options.ConsolidateAfter,
		},
		&cmd.BoolFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "enable-spottospot",
				Description: "enables the spot to spot consolidation feature",
			},
			Option: &options.EnableSpotToSpot,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "expire-after",
				Description: "time a node can live on the cluster before being deleted",
			},
			Option: &options.ExpireAfter,
		},
		&cmd.IntFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "replicas",
				Description: "number of replicas for the controller deployment",
			},
			Option: &options.Replicas,
		},
	}
	return
}
