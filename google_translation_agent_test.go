package main

import (
	"golang.org/x/text/language"
	"testing"
	"time"
)

func TestGoogleTranslationAgent_Execute(t *testing.T) {
	task := &Task{
		Id: "1",
		Agent: Agent{
			AgentId: "1",
			Name:    "Sample task",
		},
		Executor: &GoogleTranslationAgent{},
		Args:     map[string]interface{}{"source": language.English, "target": language.French},
	}

	cancel := task.Run()

	// Cancel the task after 3 seconds
	time.AfterFunc(3*time.Second, cancel)

	// Wait for a while to see the results
	time.Sleep(5 * time.Second)
}
