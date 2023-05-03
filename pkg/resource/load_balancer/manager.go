package load_balancer

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun             bool
	elbClientv1        *aws.ElasticloadbalancingClient
	elbClientv2        *aws.Elasticloadbalancingv2Client
	loadBalancerGetter *Getter
}

func (m *Manager) Init() {
	if m.elbClientv1 == nil {
		m.elbClientv1 = aws.NewElasticloadbalancingClientv1()
	}
	if m.elbClientv2 == nil {
		m.elbClientv2 = aws.NewElasticloadbalancingClientv2()
	}
	m.loadBalancerGetter = NewGetter(m.elbClientv1, m.elbClientv2)
}

func (m *Manager) Create(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) Delete(options resource.Options) (err error) {
	lbName := options.Common().Name

	elbs, err := m.loadBalancerGetter.GetLoadBalancerByName(lbName)
	if err != nil {
		return err
	}

	if len(elbs.V1) > 0 {
		err = m.elbClientv1.DeleteLoadBalancer(lbName)
	} else {
		err = m.elbClientv2.DeleteLoadBalancer(awssdk.ToString(elbs.V2[0].LoadBalancerArn))
	}

	if err != nil {
		return err
	}
	fmt.Println("Load balancer deleted successfully")

	return nil
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}
