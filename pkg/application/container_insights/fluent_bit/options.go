package fluent_bit

import (
	"github.com/awslabs/eksdemo/pkg/application"
	"github.com/awslabs/eksdemo/pkg/cmd"
)

type FluentBitOptions struct {
	*application.ApplicationOptions
	HttpServer     string
	HttpServerPort string
	ReadFromHead   string
	ReadFromTail   string
}

func newOptions() (options *FluentBitOptions, flags cmd.Flags) {
	options = &FluentBitOptions{
		HttpServer:     "On",
		HttpServerPort: "2020",
		ReadFromHead:   "Off",
		ReadFromTail:   "On",
		ApplicationOptions: &application.ApplicationOptions{
			DisableServiceAccountFlag: true,
			DisableVersionFlag:        true,
			Namespace:                 "amazon-cloudwatch",
			ServiceAccount:            "fluent-bit",
		},
	}
	// TODO: add flags for FluentBitOptions
	flags = cmd.Flags{}

	return
}
