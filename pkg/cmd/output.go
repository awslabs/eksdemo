package cmd

import "github.com/awslabs/eksdemo/pkg/printer"

type OutputFlag printer.Output

func NewOutputFlag(o *printer.Output) *OutputFlag {
	*o = printer.Table
	return (*OutputFlag)(o)
}

func (o *OutputFlag) String() string {
	return string(*o)
}

func (o *OutputFlag) Set(outputType string) error {
	out, err := printer.NewOutput(outputType)
	if err != nil {
		return err
	}
	*o = OutputFlag(out)
	return nil
}

func (o *OutputFlag) Type() string {
	return "string"
}
