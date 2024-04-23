package session

import (
	"fmt"
	"strconv"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/session-manager-plugin/src/datachannel"
	"github.com/aws/session-manager-plugin/src/log"
	"github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/resource"
	"github.com/google/uuid"
	"github.com/spf13/cobra"

	// Required for initialization
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/shellsession"
)

type Manager struct {
	DryRun bool
	Getter
}

func (m *Manager) Create(options resource.Options) error {
	sessOptions, ok := options.(*SessionOptions)
	if !ok {
		return fmt.Errorf("internal error, unable to cast options to SessionOptions")
	}

	instanceID := options.Common().Name
	params := map[string][]string{}

	if sessOptions.PortForward != 0 {
		params["portNumber"] = []string{strconv.Itoa(sessOptions.PortForward)}
		params["localPortNumber"] = []string{strconv.Itoa(sessOptions.PortForwardLocal)}
	}

	if m.DryRun {
		return m.dryRun(instanceID, sessOptions, params)
	}

	ssmClient := aws.NewSSMClient()
	out, err := ssmClient.StartSession(sessOptions.DocumentName, instanceID, params)
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

func (m *Manager) dryRun(instanceID string, options *SessionOptions, parameters map[string][]string) error {
	fmt.Println("\nSSM Session Resource Manager Dry Run:")

	fmt.Printf("SSM API Call %q with request parameters:\n", "CreateSession")
	fmt.Printf("DocumentName: %q\n", options.DocumentName)
	fmt.Printf("Target: %q\n", instanceID)
	for i, j := range parameters {
		fmt.Printf("Parameters[%q]: %q\n", i, j)
	}
	fmt.Printf("Then a websocket connection is started\n\n")

	return nil
}
