package kubernetes

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// ClientConfig is used to make it easy to get an api server client
//
// type ClientConfig interface {
// 	// RawConfig returns the merged result of all overrides
// 	RawConfig() (clientcmdapi.Config, error)
// 	// ClientConfig returns a complete client config
// 	ClientConfig() (*restclient.Config, error)
// 	// Namespace returns the namespace resulting from the merged
// 	// result of all overrides and a boolean indicating if it was
// 	// overridden
// 	Namespace() (string, bool, error)
// 	// ConfigAccess returns the rules for loading/persisting the config.
// 	ConfigAccess() ConfigAccess
// }

// Raw is clientcmdapi.Config -- represents kubeconfig

func Client(context string) (*kubernetes.Clientset, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(restConfig)
}

func ClusterURLForCurrentContext() string {
	raw, err := Kubeconfig()
	if err != nil {
		return ""
	}

	if err := clientcmdapi.MinifyConfig(raw); err != nil {
		return ""
	}

	return raw.Clusters[raw.Contexts[raw.CurrentContext].Cluster].Server
}

func DynamicClient(context string) (dynamic.Interface, error) {
	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		},
	)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	return dynamic.NewForConfig(restConfig)
}

func KubeContextForCluster(cluster *types.Cluster) (string, error) {
	raw, err := Kubeconfig()
	if err != nil {
		return "", err
	}

	found := ""

	for name, context := range raw.Contexts {
		if _, ok := raw.Clusters[context.Cluster]; ok {
			if raw.Clusters[context.Cluster].Server == aws.ToString(cluster.Endpoint) {
				found = name
				break
			}
		}
	}

	return found, nil
}

func Kubeconfig() (*clientcmdapi.Config, error) {
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	raw, err := config.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	return &raw, nil
}
