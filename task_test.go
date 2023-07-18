package main

import (
	"testing"
	"time"
)

func TestMockExecutor_Execute(t *testing.T) {
	task := &Task{
		Id:       "1",
		Name:     "Sample task",
		Executor: &MockExecutor{},
	}

	cancel := task.Run()

	// Cancel the task after 3 seconds
	time.AfterFunc(3*time.Second, cancel)

	// Wait for a while to see the results
	time.Sleep(5 * time.Second)
}
