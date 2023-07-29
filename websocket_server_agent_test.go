package main

import (
	"testing"
	"time"
)

func TestNewWebSocketServerAgent(t *testing.T) {

	spec := TaskSpec{
		Name: "WebSocketServerAgent",
		Args: map[string]interface{}{
			"address": "ws://localhost:8080",
		},
	}
	task := &Task{
		Id:       "1",
		Agent:    Agent{AgentId: "1", Name: "WebSocket Server Task"},
		Executor: NewWebSocketServerAgent(spec),
	}

	cancel := task.Run()

	// Cancel the task after 1 minute
	time.AfterFunc(1*time.Minute, cancel)

	// Wait for a while to see the results
	time.Sleep(2 * time.Minute)
}
