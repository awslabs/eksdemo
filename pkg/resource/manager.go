package resource

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Manager interface {
	Create(Options) error
	Delete(Options) error
	Init()
	SetDryRun()
	Update(Options, *cobra.Command) error
}

type CreateNotSupported struct{}

func (*CreateNotSupported) Create(_ Options) error {
	return fmt.Errorf("create not supported for this resource")
}

type DeleteNotSupported struct{}

func (*DeleteNotSupported) Delete(_ Options) error {
	return fmt.Errorf("delete not supported for this resource")
}

type UpdateNotSupported struct{}

func (*UpdateNotSupported) Update(_ Options, _ *cobra.Command) error {
	return fmt.Errorf("update not supported for this resource")
}

type DryRun struct {
	DryRun bool
}

func (m *DryRun) SetDryRun() {
	m.DryRun = true
}
