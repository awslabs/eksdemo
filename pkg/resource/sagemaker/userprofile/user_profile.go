package userprofile

import (
	"github.com/awslabs/eksdemo/pkg/cmd"
	"github.com/awslabs/eksdemo/pkg/resource"
)

func New() *resource.Resource {
	options, deleteFlags, getFlags := newOptions()

	return &resource.Resource{
		Command: cmd.Command{
			Name:        "user-profile",
			Description: "SageMaker User Profile",
			Aliases:     []string{"user-profiles", "userprofiles", "userprofile", "up"},
			Args:        []string{"USER_PROFILE_NAME"},
		},

		DeleteFlags: deleteFlags,
		GetFlags:    getFlags,

		Getter: &Getter{},

		Manager: &Manager{},

		Options: options,
	}
}
