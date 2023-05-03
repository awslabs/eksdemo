package results

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource/logs_insights/stats"
)

type ResultsPrinter struct {
	results   *cloudwatchlogs.GetQueryResultsOutput
	field     string
	logStream string
	queryId   string
	showStats bool
}

func NewPrinter(results *cloudwatchlogs.GetQueryResultsOutput, field, logStream, queryId string, showStats bool) *ResultsPrinter {
	return &ResultsPrinter{results, field, logStream, queryId, showStats}
}

func (p *ResultsPrinter) PrintTable(writer io.Writer) error {
	logLines := 0
	fieldCount := map[string]int{}
	logStreamMessageCount := map[string]int{}

	for _, result := range p.results.Results {
		logStream := "@logStream field missing"
		fieldMatch := false
		message := ""

		for _, rf := range result {
			f := aws.ToString(rf.Field)
			fieldCount[f]++

			switch f {

			case p.field:
				fieldMatch = true
				message = aws.ToString(rf.Value)
				if p.field == "@logStream" {
					logStream = aws.ToString(rf.Value)
				}

			case "@logStream":
				logStream = aws.ToString(rf.Value)
			}
		}

		if fieldMatch {
			logStreamMessageCount[logStream]++

			if p.logStream != "" && p.logStream != logStream {
				continue
			}

			logLines++
			fmt.Printf("%d | %s\n", logLines, message)
		}
	}

	if p.results.Statistics.RecordsMatched == 0 {
		fmt.Println("No results.")
	} else if logLines == 0 {
		fmt.Printf("No results for %q field.\n", p.field)
	}

	if p.showStats {
		stats.NewPrinter(p.results, p.queryId).PrintTable(os.Stdout)
	} else {
		fmt.Println("---")
	}

	if p.results.Statistics.RecordsMatched > 0 {
		if len(logStreamMessageCount) > 0 {
			logStreams := []string{}
			for k, v := range logStreamMessageCount {
				logStreams = append(logStreams, fmt.Sprintf("%s (%d)", k, v))
			}

			sort.Strings(logStreams)
			fmt.Printf("Log Stream(s):\n| %s\n", strings.Join(logStreams, "\n| "))
		}

		fields := []string{}
		for k, v := range fieldCount {
			fields = append(fields, fmt.Sprintf("%s (%d)", k, v))
		}
		fmt.Printf("Field(s): %s\n", strings.Join(fields, ", "))
	}

	return nil
}

func (p *ResultsPrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.results)
}

func (p *ResultsPrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.results)
}
