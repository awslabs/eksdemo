package ssm_session

import (
	"fmt"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Manager struct {
	DryRun bool
	Getter
}

func (m *Manager) Create(options resource.Options) error {
	instanceID := options.Common().Name

	if m.DryRun {
		return m.dryRun(instanceID)
	}

	ssmClient := aws.NewSSMClient()
	out, err := ssmClient.StartSession(instanceID)
	if err != nil {
		return err
	}

	ep, err := ssmClient.Endpoint()
	if err != nil {
		return err
	}

	ssmSession := session.Session{
		ClientId:    uuid.NewString(),
		DataChannel: &datachannel.DataChannel{},
		Endpoint:    ep.URL,
		SessionId:   awssdk.ToString(out.SessionId),
		StreamUrl:   awssdk.ToString(out.StreamUrl),
		TargetId:    instanceID,
		TokenValue:  awssdk.ToString(out.TokenValue),
	}

	return ssmSession.Execute(log.Logger(false, ssmSession.ClientId))
}

func (m *Manager) Delete(options resource.Options) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) SetDryRun() {
	m.DryRun = true
}

func (m *Manager) Update(options resource.Options, cmd *cobra.Command) error {
	return fmt.Errorf("feature not supported")
}

func (m *Manager) dryRun(instanceId string) error {
	fmt.Println("\nSSM Session Resource Manager Dry Run:")

	fmt.Printf("SSM API Call %q with request parameters:\n", "CreateSession")
	fmt.Printf("Target: %q\n", instanceId)
	fmt.Printf("Then the aws/session-manager-plugin code is used to start a websocket connection\n\n")

	return nil
}
