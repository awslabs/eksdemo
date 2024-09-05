package cmd

import "fmt"

type ArgumentAndFlagCantBeUsedTogetherError struct {
	Arg  string
	Flag string
}

func (e *ArgumentAndFlagCantBeUsedTogetherError) Error() string {
	return fmt.Sprintf("%q argument and %q flag can not be used together", e.Arg, e.Flag)
}

type FlagRequiresFlagError struct {
	Flag1 string
	Flag2 string
}

func (e *FlagRequiresFlagError) Error() string {
	return fmt.Sprintf("%q flag requires %q flag", e.Flag1, e.Flag2)
}

type MustIncludeEitherArgumentOrFlag struct {
	Arg  string
	Flag string
}

func (e *MustIncludeEitherArgumentOrFlag) Error() string {
	return fmt.Sprintf("must include either %q argument or %q flag", e.Arg, e.Flag)
}
