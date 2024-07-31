package controlPlane

import (
	"fmt"
	"os/exec"
	"os"
	"strings"
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/installer"
	"github.com/awslabs/eksdemo/pkg/template"
)

func NewApp() *application.Application {
/*
Begin mTLS Certificates
*/
/*
	externalCA := ""
	if externalCA == "" {
*/
/*
Trust anchor certificate
*/
		tlsCmd := exec.Command( "step",
			"certificate",
			"create",
			"root.linkerd.cluster.local",
			"./pkg/application/linkerd/linkerd_control_plane/ca.crt",
			"./pkg/application/linkerd/linkerd_control_plane/ca.key",
			"--profile",
			"root-ca",
			"--no-password",
			"--insecure",
			"--force" )
		stdout, err := tlsCmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		if stdout != nil {
			fmt.Println(string(stdout))
		}
		caPEMBytes, err := os.ReadFile("./pkg/application/linkerd/linkerd_control_plane/ca.crt")
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		//caPEM := string(caPEMBytes)
		caPEM := strings.Join( strings.Split( string(caPEMBytes),"\n")[:],"\n  " )

/*
Issuer certificate and key
*/
		tlsCmd = exec.Command( "step",
			"certificate",
			"create",
			"identity.linkerd.cluster.local",
			"./pkg/application/linkerd/linkerd_control_plane/issuer.crt",
			"./pkg/application/linkerd/linkerd_control_plane/issuer.key",
			"--profile",
			"intermediate-ca",
			"--not-after",
			"8760h",
			"--no-password",
			"--insecure",
			"--ca",
			"./pkg/application/linkerd/linkerd_control_plane/ca.crt",
			"--ca-key",
			"./pkg/application/linkerd/linkerd_control_plane/ca.key",
			"--force" )
		stdout, err = tlsCmd.Output()
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		if stdout != nil {
			fmt.Println(string(stdout))
		}

		crtPEMBytes, err := os.ReadFile("./pkg/application/linkerd/linkerd_control_plane/issuer.crt")
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		//crtPEM := string(crtPEMBytes)
		//fmt.Println( "[DEBUG] line 0:\n", strings.Join( strings.Split(crtPEM,"\n")[:],"") )
		crtPEM := strings.Join( strings.Split( string(crtPEMBytes),"\n")[:],"\n              " )

		keyPEMBytes, err := os.ReadFile("./pkg/application/linkerd/linkerd_control_plane/issuer.key")
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		//keyPEM := string(keyPEMBytes)
		keyPEM := strings.Join( strings.Split( string(keyPEMBytes),"\n")[:],"\n              " )
//	//}
/*
End mTLS Certificates
*/

/*
Begin Helm App
*/

	helmValues := fmt.Sprintf(valuesTemplate, caPEM, crtPEM, keyPEM)

	app := &application.Application{
		Command: cmd.Command{
			Parent:      "linkerd",
			Name:        "linkerd-control-plane",
			Description: "Linkerd Service Mesh Custom Resource Definitions",
		},

                Options: &application.ApplicationOptions{
                        Namespace: "linkerd",
                },

		Installer: &installer.HelmInstaller{
			ChartName:     "linkerd-control-plane",
			ReleaseName:   "linkerd-control-plane",
			RepositoryURL: "https://helm.linkerd.io/edge",
			ValuesTemplate: &template.TextTemplate{
				//Template: valuesTemplate,
				Template: helmValues,
			},
		},
	}
	app.Options, app.Flags = newOptions()
	return app
}

// https://github.com/linkerd/linkerd2/blob/main/charts/linkerd-control-plane/values.yaml
const valuesTemplate = `---
identityTrustAnchorsPEM: |
  %s
identity:
    scheme: linkerd.io/tls
    issuer:
        tls:
            crtPEM: |
              %s
            keyPEM: |
              %s
`
