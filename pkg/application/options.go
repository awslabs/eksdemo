package application

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/awslabs/eksdemo/pkg/resource/irsa"
	"github.com/spf13/cobra"
)

type Options interface {
	AddInstallFlags(*cobra.Command, cmd.Flags, InstallerType) cmd.Flags
	AddUninstallFlags(*cobra.Command, cmd.Flags, bool) cmd.Flags
	AssignCommonResourceOptions(*resource.Resource)
	Common() *ApplicationOptions
	KubeContext() string
	PreDependencies(Action) error
	PreInstall() error
	PostInstall(string, []*resource.Resource) error
}

type ApplicationOptions struct {
	ChartVersion string
	Version      string

	DefaultVersion
	IngressOptions

	DeleteDependencies           bool
	DisableNamespaceFlag         bool
	DisableServiceAccountFlag    bool
	DisableVersionFlag           bool
	ExposeIngressAndLoadBalancer bool
	ExposeIngressOnly            bool
	LockVersionFlag              bool
	SetValues                    []string
	UsePrevious                  bool

	Account        string
	ClusterName    string
	DryRun         bool
	Namespace      string
	Partition      string
	Region         string
	ServiceAccount string
	Cluster        *types.Cluster
	kubeContext    string
}

type Action string

const Install Action = "install"
const Uninstall Action = "uninstall"

func (o *ApplicationOptions) AddInstallFlags(cobraCmd *cobra.Command, flags cmd.Flags, it InstallerType) cmd.Flags {
	// Cluster flag has to be ordered before Version flag as it depends on the EKS cluster version
	flags = append(cmd.Flags{o.NewClusterFlag(Install), o.NewDryRunFlag()}, flags...)

	if !o.DisableVersionFlag {
		flags = append(flags, o.NewVersionFlag(), o.NewUsePreviousFlag())
	}

	if !o.DisableNamespaceFlag {
		flags = append(flags, o.NewNamespaceFlag(Install))
	}

	if !o.DisableServiceAccountFlag {
		flags = append(flags, o.NewServiceAccountFlag())
	}

	if o.ExposeIngressAndLoadBalancer || o.ExposeIngressOnly {
		o.IngressOptions = NewIngressOptions(o.ExposeIngressOnly)
		flags = append(flags, o.NewIngressFlags()...)
	}

	if it == HelmInstaller {
		flags = append(flags, o.NewChartVersionFlag(), o.NewSetFlag())
	}

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *ApplicationOptions) AddUninstallFlags(cobraCmd *cobra.Command, _ cmd.Flags, iamPolicy bool) cmd.Flags {
	commonFlags := cmd.Flags{
		o.NewClusterFlag(Uninstall),
		o.NewNamespaceFlag(Uninstall),
	}

	if iamPolicy {
		commonFlags = append(commonFlags, o.NewDeleteRoleFlag())
	}

	flags := commonFlags

	for _, f := range flags {
		f.AddFlagToCommand(cobraCmd)
	}

	return flags
}

func (o *ApplicationOptions) AssignCommonResourceOptions(res *resource.Resource) {
	if o.DryRun {
		res.SetDryRun()
	}

	r := res.Common()

	r.Account = o.Account
	r.Cluster = o.Cluster
	r.ClusterName = o.ClusterName
	r.KubeContext = o.kubeContext
	r.Namespace = o.Namespace
	r.Partition = o.Partition
	r.Region = o.Region

	// Allow for multiple IRSA for an application
	// By default will use the application service account, unless already set on the IRSA resource
	if r.ServiceAccount == "" {
		r.ServiceAccount = o.ServiceAccount
	}
}

func (o *ApplicationOptions) Common() *ApplicationOptions {
	return o
}

func (o *ApplicationOptions) IrsaAnnotation() string {
	irsaOptions := irsa.IrsaOptions{
		CommonOptions: resource.CommonOptions{
			Account:        o.Account,
			ClusterName:    o.ClusterName,
			Namespace:      o.Namespace,
			Partition:      o.Partition,
			ServiceAccount: o.ServiceAccount,
		},
	}
	return irsaOptions.IrsaAnnotation()
}

func (o *ApplicationOptions) IrsaAnnotationFor(serviceAccount string) string {
	irsaOptions := irsa.IrsaOptions{
		CommonOptions: resource.CommonOptions{
			Account:        o.Account,
			ClusterName:    o.ClusterName,
			Namespace:      o.Namespace,
			Partition:      o.Partition,
			ServiceAccount: serviceAccount,
		},
	}
	return irsaOptions.IrsaAnnotation()
}

func (o *ApplicationOptions) KubeContext() string {
	return o.kubeContext
}

func (o *ApplicationOptions) PreDependencies(Action) error {
	return nil
}

func (o *ApplicationOptions) PreInstall() error {
	return nil
}

func (o *ApplicationOptions) PostInstall(name string, postInstallRes []*resource.Resource) error {
	if res := o.IngressOptions.PostInstallResources(name); res != nil {
		postInstallRes = append(postInstallRes, res...)
	}

	if len(postInstallRes) > 0 {
		fmt.Printf("Creating %d post-install resources for %s\n", len(postInstallRes), name)
	}

	for _, res := range postInstallRes {
		fmt.Printf("Creating post-install resource: %s\n", res.Common().Name)

		o.AssignCommonResourceOptions(res)

		if err := res.Create(); err != nil {
			return err
		}
	}
	return nil
}
