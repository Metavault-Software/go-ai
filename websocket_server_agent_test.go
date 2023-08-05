package main

import (
	"testing"
	"time"
)

func TestNewWebSocketServerAgent(t *testing.T) {
	spec := Agent{
		Name: "WebSocketServerAgent",
		Args: map[string]interface{}{
			"address": "ws://localhost:8080",
		},
	}
	task := &Task{
		Id:    "1",
		Agent: spec,
	}
	err := task.ExecutorFromExecutorType()
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	cancel := task.Run()

	// Cancel the task after 1 minute
	time.AfterFunc(1*time.Minute, cancel)

	// Wait for a while to see the results
	time.Sleep(2 * time.Minute)
}
