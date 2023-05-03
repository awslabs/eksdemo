package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type Flag interface {
	AddFlagToCommand(*cobra.Command)
	ValidateFlag(cmd *cobra.Command, args []string) error
	GetName() string
}

type Flags []Flag

type CommandFlag struct {
	Name        string
	Description string
	Shorthand   string
	Required    bool
	Validate    func(cmd *cobra.Command, args []string) error
}

const required = " (required)"

type BoolFlag struct {
	CommandFlag
	Option *bool
}

type DurationFlag struct {
	CommandFlag
	Option *time.Duration
}

type IntFlag struct {
	CommandFlag
	Option *int
}

type StringFlag struct {
	CommandFlag
	Choices []string
	Option  *string
}

type StringSliceFlag struct {
	CommandFlag
	Choices []string
	Option  *[]string
}

// Command flag methods
func (f *CommandFlag) GetName() string {
	return f.Name
}

// Boolean flag methods
func (f *BoolFlag) AddFlagToCommand(cmd *cobra.Command) {
	if f.Shorthand != "" {
		cmd.Flags().BoolVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	} else {
		cmd.Flags().BoolVar(f.Option, f.Name, *f.Option, f.Description)
	}
}

func (f *BoolFlag) ValidateFlag(cmd *cobra.Command, args []string) error {
	if f.Validate == nil {
		return nil
	}
	return f.Validate(cmd, args)
}

// Duration flag methods
func (f *DurationFlag) AddFlagToCommand(cmd *cobra.Command) {
	if f.Required {
		f.Description += required
	}

	if f.Shorthand != "" {
		cmd.Flags().DurationVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	} else {
		cmd.Flags().DurationVar(f.Option, f.Name, *f.Option, f.Description)
	}

	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *DurationFlag) ValidateFlag(cmd *cobra.Command, args []string) error {
	if f.Validate == nil {
		return nil
	}
	return f.Validate(cmd, args)
}

// Int flag methods
func (f *IntFlag) AddFlagToCommand(cmd *cobra.Command) {
	if f.Required {
		f.Description += required
	}

	if f.Shorthand != "" {
		cmd.Flags().IntVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	} else {
		cmd.Flags().IntVar(f.Option, f.Name, *f.Option, f.Description)
	}

	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *IntFlag) ValidateFlag(cmd *cobra.Command, args []string) error {
	if f.Validate == nil {
		return nil
	}
	return f.Validate(cmd, args)
}

// String flag methods
func (f *StringFlag) AddFlagToCommand(cmd *cobra.Command) {
	if f.Required {
		f.Description += required
	}

	if f.Shorthand != "" {
		cmd.Flags().StringVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)
	} else {
		cmd.Flags().StringVar(f.Option, f.Name, *f.Option, f.Description)
	}

	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *StringFlag) ValidateFlag(cmd *cobra.Command, args []string) error {
	if len(f.Choices) > 0 {
		found := false

		for _, choice := range f.Choices {
			if strings.EqualFold(choice, *f.Option) {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("--%s must be one of: %s", f.Name, strings.Join(f.Choices, ", "))
		}
	}

	if f.Validate != nil {
		return f.Validate(cmd, args)
	}

	return nil
}

// StringSlice flag methods
func (f *StringSliceFlag) AddFlagToCommand(cmd *cobra.Command) {
	if f.Required {
		f.Description += required
	}

	cmd.Flags().StringSliceVarP(f.Option, f.Name, f.Shorthand, *f.Option, f.Description)

	if f.Required {
		cmd.MarkFlagRequired(f.Name)
	}
}

func (f *StringSliceFlag) ValidateFlag(cmd *cobra.Command, args []string) error {
	if len(f.Choices) > 0 {

		for _, flag := range *f.Option {
			found := false

			for _, choice := range f.Choices {
				if strings.EqualFold(choice, flag) {
					found = true
				}
			}

			if !found {
				return fmt.Errorf("--%s can only contain: %s", f.Name, strings.Join(f.Choices, ", "))
			}
		}
	}

	if f.Validate != nil {
		return f.Validate(cmd, args)
	}

	return nil
}

// Flags (list of flags) methods
func (f Flags) ValidateFlags(cmd *cobra.Command, args []string) error {
	for _, flag := range f {
		if err := flag.ValidateFlag(cmd, args); err != nil {
			return err
		}
	}
	return nil
}

func (f Flags) Remove(name string) Flags {
	for i, flag := range f {
		if flag.GetName() == name {
			f[i] = f[len(f)-1]
			f = f[:len(f)-1]
			break
		}
	}
	return f
}
