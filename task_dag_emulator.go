package main

import "fmt"

type TaskDagEmulator struct {
	TaskDag
}

func (e *TaskDagEmulator) Generate100RandomTasks() []Task {
	var tasks []Task
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("task-%d", i)
		task := Task{
			Id:   id,
			Name: id,
			Executor: NewEmulateOpenAIAgent(TaskSpec{
				ID:       id,
				Name:     id,
				Executor: "OpenAIAgent",
				Args: map[string]interface{}{
					"messages": []interface{}{
						map[string]interface{}{
							"role":    "user",
							"content": "Write a Go program to illustrate Go's powerful concurrency model.",
						},
					},
				},
			}),
			Done: make(chan bool),
		}
		tasks = append(tasks, task)
	}
	return tasks
}
