package iam_oidc

import (
	"fmt"

	"github.com/awslabs/eksdemo/pkg/resource"
)

type IamOidcOptions struct {
	resource.CommonOptions
}

func newOptions() (options *IamOidcOptions) {
	options = &IamOidcOptions{
		CommonOptions: resource.CommonOptions{
			ClusterFlagOptional: true,
		},
	}
	return
}

func (o *IamOidcOptions) Validate(args []string) error {
	if o.ClusterName != "" && len(args) > 0 {
		return fmt.Errorf("%q flag cannot be used with URL argument", "--cluster")
	}
	return nil
}
