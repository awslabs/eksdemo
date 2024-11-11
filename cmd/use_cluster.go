package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/kubernetes"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var useClusterCmd = &cobra.Command{
	Use:     "use-cluster CLUSTER_NAME",
	Short:   "Update kubeconfig to use the EKS cluster specified",
	Aliases: []string{"use"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterName := args[0]
		cluster, err := aws.NewEKSClient().DescribeCluster(clusterName)
		if err != nil {
			return aws.FormatErrorAsMessageOnly(err)
		}
		if cluster.Status != types.ClusterStatusActive {
			return fmt.Errorf("cluster %q is not active, cluster status is: %s", clusterName, cluster.Status)
		}

		if err := flags.ValidateFlags(cmd, args); err != nil {
			return err
		}
		cmd.SilenceUsage = true

		configAccess := clientcmd.NewDefaultPathOptions()
		configFileName := configAccess.GetDefaultFilename()

		existingConfig, err := configAccess.GetStartingConfig()
		if err != nil {
			return fmt.Errorf("failed to read existing kubeconfig %q: %w", configFileName, err)
		}

		// If a context exists for the cluster, update kubeconfig to use that context
		contextName := contextForCluster(existingConfig, cluster)
		if contextName != "" {
			existingConfig.CurrentContext = contextName
			err = clientcmd.ModifyConfig(configAccess, *existingConfig, false)
			if err != nil {
				return fmt.Errorf("failed to update kubeconfig: %w", err)
			}
			fmt.Printf("Updated Kubeconfig %q to use context %q\n", configFileName, contextName)
			return nil
		}

		// If a context doesn't exist for the cluster, create a new one in the existing Kubeconfig
		certAuthData, err := base64.StdEncoding.DecodeString(awssdk.ToString(cluster.CertificateAuthority.Data))
		if err != nil {
			return fmt.Errorf("failed decoding cluster certificate auth data: %w", err)
		}

		// Consider options in the future to configure profile and role
		profile := ""
		roleARN := ""
		newConfig := buildConfig(cluster, certAuthData, profile, roleARN)

		err = writeConfig(configAccess, existingConfig, newConfig)
		if err != nil {
			return fmt.Errorf("failed writing kubeconfig %q: %w", configFileName, err)
		}

		fmt.Printf("Added context %q to Kubeconfig %q\n", newConfig.CurrentContext, configFileName)

		return nil
	},
}

var flags = cmd.Flags{
	&cmd.StringFlag{
		CommandFlag: cmd.CommandFlag{
			Name:        "kubeconfig",
			Description: "path to write kubeconfig",
			Validate: func(_ *cobra.Command, _ []string) error {
				// Set the KUBECONFIG environment variable to configure client-go
				_ = os.Setenv("KUBECONFIG", kubeconfig)
				return nil
			},
		},
		Option: &kubeconfig,
	},
}

var kubeconfig string

func init() {
	kubeconfig = kubernetes.KubeconfigDefaultPath()
	for _, f := range flags {
		f.AddFlagToCommand(useClusterCmd)
	}
	rootCmd.AddCommand(useClusterCmd)
}

func buildConfig(cluster *types.Cluster, certAuthData []byte, profile, roleARN string) *clientcmdapi.Config {
	clusterArn := awssdk.ToString(cluster.Arn)
	contextName := clusterArn

	args := []string{
		"eks", "get-token",
		"--cluster-name", awssdk.ToString(cluster.Name),
		"--region", aws.Region(),
		"--output", "json",
	}

	if roleARN != "" {
		args = append(args, "--role-arn", roleARN)
	}

	execConfig := &clientcmdapi.ExecConfig{
		APIVersion: "client.authentication.k8s.io/v1beta1",
		Args:       args,
		Command:    "aws",
		Env: []clientcmdapi.ExecEnvVar{
			{
				Name:  "AWS_STS_REGIONAL_ENDPOINTS",
				Value: "regional",
			},
		},
		ProvideClusterInfo: false,
	}

	if profile != "" {
		execConfig.Env = append(execConfig.Env, clientcmdapi.ExecEnvVar{
			Name:  "AWS_PROFILE",
			Value: profile,
		})
	}

	return &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterArn: {
				Server:                   awssdk.ToString(cluster.Endpoint),
				CertificateAuthorityData: certAuthData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterArn,
				AuthInfo: contextName,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			contextName: {
				Exec: execConfig,
			},
		},
		CurrentContext: contextName,
	}
}

func contextForCluster(config *clientcmdapi.Config, cluster *types.Cluster) string {
	found := ""
	for name, context := range config.Contexts {
		if _, ok := config.Clusters[context.Cluster]; ok {
			if config.Clusters[context.Cluster].Server == awssdk.ToString(cluster.Endpoint) {
				found = name
				break
			}
		}
	}

	return found
}

func merge(existing *clientcmdapi.Config, tomerge *clientcmdapi.Config) *clientcmdapi.Config {
	for k, v := range tomerge.Clusters {
		existing.Clusters[k] = v
	}
	for k, v := range tomerge.AuthInfos {
		existing.AuthInfos[k] = v
	}
	for k, v := range tomerge.Contexts {
		existing.Contexts[k] = v
	}

	return existing
}

func writeConfig(configAccess clientcmd.ConfigAccess, existingConfig, newConfig *clientcmdapi.Config) error {
	merged := merge(existingConfig, newConfig)
	merged.CurrentContext = newConfig.CurrentContext

	return clientcmd.ModifyConfig(configAccess, *merged, true)
}
