package cmd

import "fmt"

type ArgumentAndFlagCantBeUsedTogetherError struct {
	Arg  string
	Flag string
}

func (e *ArgumentAndFlagCantBeUsedTogetherError) Error() string {
	return fmt.Sprintf("%q argument and %q flag can not be used together", e.Arg, e.Flag)
}
