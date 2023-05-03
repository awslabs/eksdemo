package resource

import "github.com/awslabs/eksdemo/pkg/printer"

type Getter interface {
	Get(string, printer.Output, Options) error
	Init()
}

type EmptyInit struct{}

func (i *EmptyInit) Init() {}
