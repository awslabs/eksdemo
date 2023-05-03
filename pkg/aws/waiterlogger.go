package aws

import (
	"context"
	"fmt"

	"github.com/aws/smithy-go/middleware"
)

// WaiterLogger is a Custom Logger middleware used by the waiter to log an attempt with a "."
type WaiterLogger struct {
	// Attempt is the current attempt to be logged
	Attempt int64
}

// ID representing the Logger middleware
func (*WaiterLogger) ID() string {
	return "CustomWaiterLogger"
}

// HandleInitialize performs handling of request in initialize stack step
func (m *WaiterLogger) HandleInitialize(ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	fmt.Printf(".")

	return next.HandleInitialize(ctx, in)
}

// AddLogger is a helper util to add waiter logger after `SetLogger` middleware in
func (m WaiterLogger) AddLogger(stack *middleware.Stack) error {
	return stack.Initialize.Insert(&m, "SetLogger", middleware.After)
}
