package sqs_queue

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awslabs/eksdemo/pkg/aws"
	"github.com/awslabs/eksdemo/pkg/printer"
	"github.com/awslabs/eksdemo/pkg/resource"
)

type Queue struct {
	Attributes map[string]string
	Url        string
}

type Getter struct {
	sqsClient *aws.SQSClient
}

func NewGetter(sqsClient *aws.SQSClient) *Getter {
	return &Getter{sqsClient}
}

func (g *Getter) Init() {
	if g.sqsClient == nil {
		g.sqsClient = aws.NewSQSClient()
	}
}

func (g *Getter) Get(name string, output printer.Output, options resource.Options) error {
	var queue Queue
	var queues []Queue
	var err error

	if name != "" {
		queue, err = g.GetQueueByName(name)
		queues = []Queue{queue}
	} else {
		queues, err = g.GetAllQueues()
	}

	if err != nil {
		return err
	}

	return output.Print(os.Stdout, NewPrinter(queues))
}

func (g *Getter) GetAllQueues() ([]Queue, error) {
	queuesUrls, err := g.sqsClient.ListQueues()
	if err != nil {
		return nil, err
	}

	queues := make([]Queue, 0, len(queuesUrls))

	for _, url := range queuesUrls {
		attr, err := g.sqsClient.GetQueueAttributes(url)
		if err != nil {
			return nil, err
		}
		queues = append(queues, Queue{attr, url})
	}

	return queues, nil
}

func (g *Getter) GetQueueByName(name string) (Queue, error) {
	queueUrl, err := g.sqsClient.GetQueueUrl(name)
	var qdne *types.QueueDoesNotExist
	if err != nil && errors.As(err, &qdne) {
		return Queue{}, resource.NotFoundError(fmt.Sprintf("sqs-queue %q not found", name))
	}

	attr, err := g.sqsClient.GetQueueAttributes(queueUrl)
	if err != nil {
		return Queue{}, err
	}

	return Queue{attr, queueUrl}, nil
}
