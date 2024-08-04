package controlplane

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type AppOptions struct {
	application.ApplicationOptions

	trustAnchor string
	issuerCert  string
	issuerKey   string

	CAPEM  string
	CRTPEM string
	KEYPEM string
}

func newOptions() (options *AppOptions, flags cmd.Flags) {
	options = &AppOptions{
		ApplicationOptions: application.ApplicationOptions{
			DefaultVersion: &application.LatestPrevious{
				LatestChart:   "2024.7.3",
				PreviousChart: "2024.7.2",
				Latest:        "edge-24.7.3",
				Previous:      "edge-24.7.2",
			},
			Namespace: "linkerd",
		},
	}

	flags = cmd.Flags{
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "trust-anchor",
				Description: "Path to TLS Certificate to use as the Trust Anchor",
			},
			Option: &options.trustAnchor,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "issuer-cert",
				Description: "Path to TLS Certificate to use as the Issuer",
			},
			Option: &options.issuerCert,
		},
		&cmd.StringFlag{
			CommandFlag: cmd.CommandFlag{
				Name:        "issuer-key",
				Description: "Path to TLS Key to use as the Issuer",
			},
			Option: &options.issuerKey,
		},
	}
	return
}

func (options *AppOptions) PreInstall() error {
	/*
	   Begin mTLS Certificates
	*/

	/*
	   Trust anchor certificate
	*/
	tlsCmd := exec.Command("step",
		"certificate",
		"create",
		"root.linkerd.cluster.local",
		"./pkg/application/linkerd/controlplane/ca.crt",
		"./pkg/application/linkerd/controlplane/ca.key",
		"--profile",
		"root-ca",
		"--no-password",
		"--insecure",
		"--force")
	stdout, err := tlsCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if stdout != nil {
		fmt.Println(string(stdout))
	}
	caPEMBytes, err := os.ReadFile("./pkg/application/linkerd/controlplane/ca.crt")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	caPEM := strings.Join(strings.Split(string(caPEMBytes), "\n")[:], "\n  ")

	/*
	   Issuer certificate and key
	*/
	tlsCmd = exec.Command("step",
		"certificate",
		"create",
		"identity.linkerd.cluster.local",
		"./pkg/application/linkerd/controlplane/issuer.crt",
		"./pkg/application/linkerd/controlplane/issuer.key",
		"--profile",
		"intermediate-ca",
		"--not-after",
		"8760h",
		"--no-password",
		"--insecure",
		"--ca",
		"./pkg/application/linkerd/controlplane/ca.crt",
		"--ca-key",
		"./pkg/application/linkerd/controlplane/ca.key",
		"--force")
	stdout, err = tlsCmd.Output()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if stdout != nil {
		fmt.Println(string(stdout))
	}

	crtPEMBytes, err := os.ReadFile("./pkg/application/linkerd/controlplane/issuer.crt")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	crtPEM := strings.Join(strings.Split(string(crtPEMBytes), "\n")[:], "\n              ")

	keyPEMBytes, err := os.ReadFile("./pkg/application/linkerd/controlplane/issuer.key")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	keyPEM := strings.Join(strings.Split(string(keyPEMBytes), "\n")[:], "\n              ")
	/*
	   End mTLS Certificates
	*/

	options.CAPEM = caPEM
	options.CRTPEM = crtPEM
	options.KEYPEM = keyPEM
	return nil
}
