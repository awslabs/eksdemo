package sqs_queue

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/hako/durafmt"
)

type SqsQueuePrinter struct {
	queues []Queue
}

func NewPrinter(queues []Queue) *SqsQueuePrinter {
	return &SqsQueuePrinter{queues}
}

func (p *SqsQueuePrinter) PrintTable(writer io.Writer) error {
	table := printer.NewTablePrinter()
	table.SetHeader([]string{"Age", "Name", "Type", "Messages", "In Flight"})

	for _, q := range p.queues {
		age := "failed to parse"
		num, err := strconv.ParseInt(q.Attributes[string(types.QueueAttributeNameCreatedTimestamp)], 10, 64)
		if err == nil {
			age = durafmt.ParseShort(time.Since(time.Unix(num, 0))).String()
		}

		name := "failed to parse"
		queueArn, err := arn.Parse(q.Attributes[string(types.QueueAttributeNameQueueArn)])
		if err == nil {
			name = queueArn.Resource
		}

		queueType := "Standard"
		if strings.HasSuffix(name, ".fifo") {
			queueType = "FIFO"
		}

		table.AppendRow([]string{
			age,
			name,
			queueType,
			q.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessages)],
			q.Attributes[string(types.QueueAttributeNameApproximateNumberOfMessagesNotVisible)],
		})
	}

	table.Print(writer)

	return nil
}

func (p *SqsQueuePrinter) PrintJSON(writer io.Writer) error {
	return printer.EncodeJSON(writer, p.queues)
}

func (p *SqsQueuePrinter) PrintYAML(writer io.Writer) error {
	return printer.EncodeYAML(writer, p.queues)
}
