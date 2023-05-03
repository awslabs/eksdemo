package nodegroup

import (
	"fmt"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	Eksctl          resource.Manager
	nodegroupGetter *Getter
	eksClient       *aws.EKSClient
}

func (m *Manager) Init() {
	if m.eksClient == nil {
		m.eksClient = aws.NewEKSClient()
	}
	m.nodegroupGetter = NewGetter(m.eksClient)
}

func (m *Manager) Create(options resource.Options) error {
	return m.Eksctl.Create(options)
}

func (m *Manager) Delete(options resource.Options) error {
	return m.Eksctl.Delete(options)
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	ngOptions, ok := options.(*NodegroupOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to NodegroupOptions")
	}

	cluster := options.Common().ClusterName
	nodegroup := ngOptions.NodegroupName

	ng, err := m.nodegroupGetter.GetNodeGroupByName(nodegroup, cluster)
	if err != nil {
		return err
	}

	unsetFlags := 0
	update := ""

	if cmd.Flags().Changed("nodes") {
		update += fmt.Sprintf("%d Nodes", ngOptions.UpdateDesired)
	} else {
		ngOptions.UpdateDesired = int(awssdk.ToInt32(ng.ScalingConfig.DesiredSize))
		unsetFlags++
	}

	if cmd.Flags().Changed("min") {
		if len(update) > 0 {
			update += ", "
		}
		update += fmt.Sprintf("%d Min", ngOptions.MinSize)
	} else {
		ngOptions.UpdateMin = int(awssdk.ToInt32(ng.ScalingConfig.MinSize))
		unsetFlags++
	}

	if cmd.Flags().Changed("max") {
		if len(update) > 0 {
			update += ", "
		}
		update += fmt.Sprintf("%d Max", ngOptions.MaxSize)
	} else {
		ngOptions.UpdateMax = int(awssdk.ToInt32(ng.ScalingConfig.MaxSize))
		unsetFlags++
	}

	if unsetFlags == 3 {
		return fmt.Errorf("at least one flag %s is required", strings.Join([]string{"\"nodes\"", "\"min\"", "\"max\""}, ", "))
	}

	fmt.Printf("Updating nodegroup with %s...", update)

	err = m.eksClient.UpdateNodegroupConfig(cluster, nodegroup, ngOptions.UpdateDesired, ngOptions.UpdateMin, ngOptions.UpdateMax)
	if err != nil {
		return err
	}
	fmt.Println("done")

	return nil
}

func (m *Manager) SetDryRun() {
	m.Eksctl.SetDryRun()
}
