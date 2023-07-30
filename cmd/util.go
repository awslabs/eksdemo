package cmd

import (
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/util"
	"github.com/spf13/cobra"
)

var clusterName string
var namespace string

var utilCmd = &cobra.Command{
	Use:     "utils",
	Short:   "Utility commands",
	Aliases: []string{"util"},
}

var enablePrefixAssignmentCmd = &cobra.Command{
	Use:     "enable-prefix-assignment",
	Short:   "Enable Prefix Assignmenr with VPC CNI",
	Aliases: []string{"enable-prefix"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.NewEKSClient().DescribeCluster(clusterName)
		if err != nil {
			return err
		}
		return util.EnablePrefixAssignment(cluster)
	},
}

var enableSecurityGroupsForPodsCmd = &cobra.Command{
	Use:     "enable-sg-for-pods",
	Short:   "Enable Security Groups for Pods with VPC CNI",
	Aliases: []string{"enable-sgpods"},
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.NewEKSClient().DescribeCluster(clusterName)
		if err != nil {
			return err
		}
		return util.EnableSecurityGroupsForPods(cluster)
	},
}

var serviceAccountToken = &cobra.Command{
	Use:     "service-account-token SERVICE_ACCOUNT",
	Short:   "Print out the token for the service account",
	Aliases: []string{"sa-token"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cluster, err := aws.NewEKSClient().DescribeCluster(clusterName)
		if err != nil {
			return err
		}
		return util.ServiceAccountToken(cluster, namespace, args[0])
	},
}

var tagSubnetsCmd = &cobra.Command{
	Use:     "tag-subnets",
	Short:   "Add kubernetes.io/cluster/<cluster-name> tag to private VPC subnets",
	Aliases: []string{"tag-subnet"},
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := aws.NewEKSClient().DescribeCluster(clusterName)
		if err != nil {
			return err
		}

		return util.TagSubnets(clusterName)
	},
}

func init() {
	rootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(tagSubnetsCmd)
	utilCmd.AddCommand(enablePrefixAssignmentCmd)
	utilCmd.AddCommand(enableSecurityGroupsForPodsCmd)
	utilCmd.AddCommand(serviceAccountToken)

	tagSubnetsCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	tagSubnetsCmd.MarkFlagRequired("cluster")

	enablePrefixAssignmentCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	enablePrefixAssignmentCmd.MarkFlagRequired("cluster")

	enableSecurityGroupsForPodsCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	enableSecurityGroupsForPodsCmd.MarkFlagRequired("cluster")

	serviceAccountToken.Flags().StringVarP(&clusterName, "cluster", "c", "", "cluster (required)")
	serviceAccountToken.MarkFlagRequired("cluster")
	serviceAccountToken.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace (required)")
	serviceAccountToken.MarkFlagRequired("namespace")
}
