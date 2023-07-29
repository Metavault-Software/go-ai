package main

import (
	"fmt"
	"time"
)

type TaskDagEmulator struct {
	TaskGraph
}

func (e *TaskDagEmulator) Generate100RandomTasks() []Task {
	var tasks []Task
	for i := 0; i < 10; i++ {
		id := fmt.Sprintf("task-%d", i)
		task := Task{
			Id: id,
			Agent: Agent{
				Name: id,
			},
			Status: make(chan TaskStatus),
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

	// write status to Status channel
	for _, task := range tasks {
		go func(task Task) {
			for {
				task.Status <- Running
				fmt.Printf("Task %s status: %v\n", task.Id, task.Status)
				time.Sleep(1 * time.Second)
			}
		}(task)
	}

	go func() {
		// after 10 secs close all tasks
		<-time.After(10 * time.Second)
		for _, task := range tasks {
			close(task.Done)
		}
	}()

	return tasks
}
