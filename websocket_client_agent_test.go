package main

import (
	"testing"
	"time"
)

func TestWebSocketClientAgent_Execute(t *testing.T) {
	task := &Task{
		Id:       "1",
		Agent:    Agent{AgentId: "1", Name: "WebSocket Client Task"},
		Executor: &WebSocketClientAgent{Addr: "ws://localhost:8080/ws"},
	}

	cancel := task.Run()

	// Cancel the task after 1 minute
	time.AfterFunc(1*time.Minute, cancel)

	// Wait for a while to see the results
	time.Sleep(2 * time.Minute)
}
