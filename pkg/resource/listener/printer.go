package listener

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
	"github.com/awslabs/eksdemo/pkg/printer"
)

type ListenerPrinter struct {
	listeners []types.Listener
}

func NewPrinter(listeners []types.Listener) *ListenerPrinter {
	return &ListenerPrinter{listeners}
}

func (p *ListenerPrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Id", "Prot:Port", "Default Certificate Id", "Default Action"})

	resourceId := regexp.MustCompile(`[^:/]*$`)

	for _, l := range p.listeners {
		// DescribeListeners API documentation states that only the default certificate is included
		// https://docs.aws.amazon.com/elasticloadbalancing/latest/APIReference/API_DescribeListeners.html
		defaultCert := "-"
		if len(l.Certificates) > 0 {
			defaultCert = resourceId.FindString(aws.ToString(l.Certificates[0].CertificateArn))
		}

		table.AppendRow([]string{
			resourceId.FindString(aws.ToString(l.ListenerArn)),
			string(l.Protocol) + ":" + strconv.Itoa((int(aws.ToInt32(l.Port)))),
			defaultCert,
			strings.Join(PrintActions(l.DefaultActions), "\n"),
		})
	}

	table.SeparateRows()
	table.Print(writer)

	return nil
}

func (p *ListenerPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.listeners)
}

func (p *ListenerPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.listeners)
}

func PrintActions(elbActions []types.Action) (actions []string) {
	for _, a := range elbActions {
		switch {
		case a.AuthenticateCognitoConfig != nil || a.AuthenticateOidcConfig != nil:
			actions = append(actions, "TODO: authenticate action")

		case a.FixedResponseConfig != nil:
			actions = append(actions, "return fixed response "+aws.ToString(a.FixedResponseConfig.StatusCode))

		case a.ForwardConfig != nil:
			tgNames := []string{}
			for _, tg := range a.ForwardConfig.TargetGroups {
				tgNames = append(tgNames, tgName(aws.ToString(tg.TargetGroupArn)))
			}
			actions = append(actions, "forward to "+strings.Join(tgNames, "\n"))

		case a.RedirectConfig != nil:
			prot := aws.ToString(a.RedirectConfig.Protocol)
			host := aws.ToString(a.RedirectConfig.Host)
			port := aws.ToString(a.RedirectConfig.Port)
			path := aws.ToString(a.RedirectConfig.Path)
			query := aws.ToString(a.RedirectConfig.Query)
			actions = append(actions, "redirect to "+prot+"://"+host+":"+port+path+"?"+query)
		}
	}
	return
}

func tgName(tgArn string) string {
	parts := strings.Split(tgArn, "/")
	if len(parts) < 2 {
		return "failed to parse arn"
	}
	return parts[1]
}
