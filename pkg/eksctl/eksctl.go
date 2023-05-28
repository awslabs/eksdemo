package eksctl

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/hashicorp/go-version"
)

const minVersion = "0.143.0"

func GetClusterName(cluster string) string {
	return fmt.Sprintf("%s.%s.eksctl.io", cluster, aws.Region())
}

func TagNamePrefix(clusterName string) string {
	return fmt.Sprintf("eksctl-%s-cluster/", clusterName)
}

func CheckVersion() error {
	errmsg := fmt.Errorf("eksdemo requires eksctl version %s or later", minVersion)

	eksctlVersion, err := exec.Command("eksctl", "version").Output()
	if err != nil {
		return errmsg
	}

	v, err := version.NewVersion(strings.TrimSpace(string(eksctlVersion)))
	if err != nil {
		fmt.Printf("Warning: unable to parse eksctl version :%s\n", err)
		return nil
	}

	if v.LessThan(version.Must(version.NewVersion(minVersion))) {
		return errmsg
	}

	return nil
}
