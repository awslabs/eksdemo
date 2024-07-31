package key

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func NewResource() *resource.Resource {
	options, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "kms-key",
			Description: "KMS Key",
			Aliases:     []string{"kms-keys", "kmskeys", "kmskey", "kms"},
			Args:        []string{"ALIAS"},
			CreateArgs:  []string{"ALIAS"},
		},

		GetFlags: getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
