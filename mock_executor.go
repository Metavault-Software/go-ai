package main

import (
	"context"
	"fmt"
	"time"
)

// MockExecutor simulates a task that takes time to run and can be cancelled.
type MockExecutor struct{}

func (e *MockExecutor) Execute(ctx context.Context, task *Task) error {
	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			fmt.Printf("Task %s is running\n", task.Name)
			// Simulate task processing
		}
	}
	return nil
}
