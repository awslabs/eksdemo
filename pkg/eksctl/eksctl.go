package eksctl

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/hashicorp/go-version"
)

const minVersion = "0.160.0"

func GetClusterName(cluster string) string {
	return fmt.Sprintf("%s.%s.eksctl.io", cluster, aws.Region())
}

func TagNamePrefix(clusterName string) string {
	return fmt.Sprintf("eksctl-%s-cluster/", clusterName)
}

func CheckVersion() error {
	errmsg := fmt.Errorf("eksdemo requires eksctl version %s or later", minVersion)

	eksctlVersionRaw, err := exec.Command("eksctl", "version").Output()
	if err != nil {
		return errmsg
	}

	// Sometimes homebrew installs eksctl from homebrew-core
	// and `eksctl version` results in something like "0.144.0-dev+92e3cd383.2023-06-09T12:46:50Z"
	// This removes the -dev part if it exists
	eksctlVersion := strings.Split(strings.TrimSpace(string(eksctlVersionRaw)), "-")[0]

	v, err := version.NewVersion(eksctlVersion)
	if err != nil {
		fmt.Printf("Warning: unable to parse eksctl version: %s\n", err.Error())
		return nil
	}

	if v.LessThan(version.Must(version.NewVersion(minVersion))) {
		return errmsg
	}

	return nil
}
