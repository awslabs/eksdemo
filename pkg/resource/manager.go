package resource

import (
	"github.com/spf13/cobra"
)

type Manager interface {
	Create(options Options) error
	Delete(options Options) error
	Init()
	SetDryRun()
	Update(options Options, cmd *cobra.Command) error
}
