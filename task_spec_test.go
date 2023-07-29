package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestSpec(t *testing.T) {
	data, err := ioutil.ReadFile("tasks.json")
	if err != nil {
		panic(err)
	}

	var specs []TaskSpec
	err = json.Unmarshal(data, &specs)
	if err != nil {
		panic(err)
	}

	tasks := make([]*Task, len(specs))
	for i, spec := range specs {
		var executor Executor
		switch spec.Executor {
		case "OpenAIAgent":
			executor = &OpenAIAgent{}
		case "GoogleTranslationAgent":
			executor = &GoogleTranslationAgent{}
		default:
			panic(fmt.Errorf("unknown executor: %s", spec.Executor))
		}
		tasks[i] = &Task{
			Id: spec.ID,
			Agent: Agent{
				AgentId: spec.ID,
				Name:    spec.Name,
			},
			Executor: executor,
			Args:     spec.Args,
		}
	}

	// Now you can run the tasks...
}
