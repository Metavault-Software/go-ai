package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type TaskDag struct {
	Tasks map[string]Task
	Edges map[string][]Task
}

func (d *TaskDag) RunConcurrently(started, result chan Task) {
	visited := &sync.Map{}

	var visit func(task Task)
	visit = func(task Task) {
		if _, ok := visited.LoadOrStore(task.Id, true); !ok {
			for _, dep := range d.Edges[task.Id] {
				visit(dep)
			}

			// Process task: here it's a random delay of up to 5 seconds
			go func(task Task) {
				// make sure all dependencies are done
				for _, task := range d.Edges[task.Id] {
					<-task.Done
				}
				started <- task
				// simulate work
				time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
				// signal result
				result <- task
				// close done channel
				close(task.Done)
			}(task)
		}
	}

	go func() {
		for _, task := range d.Tasks {
			visit(task)
		}
		// Wait for all tasks to finish
		for _, task := range d.Tasks {
			<-task.Done
		}
		close(result)
		close(started)
	}()
}

func (d *TaskDag) FromTaskDagSpec(taskDagSpec TaskDagSpec) error {
	// Reset tasks and edges in the current TaskDag
	d.Tasks = make(map[string]Task)
	d.Edges = make(map[string][]Task)

	// Fill the Tasks map with Task ID as the key and Task object as the value
	for _, spec := range taskDagSpec.Tasks {
		var executor Executor
		switch spec.Executor {
		case "OpenAIAgent":
			executor = NewOpenAIAgent(spec)
		case "GoogleTranslationAgent":
			executor = NewGoogleTranslationAgent(spec)
		case "WebSocketClientAgent":
			executor = NewWebSocketClientAgent(spec)
		case "OpenAIGenerativeAgent":
			executor = NewOpenAIGenerativeAgent(spec)
		case "WebSocketServerAgent":
			executor = NewWebSocketServerAgent(spec)
		case "DockerAgent":
			executor = NewDockerAgent(spec)

		case "LocalFileSystemAgent":
			executor = NewLocalFileSystemAgent(spec)
		default:
			return fmt.Errorf("unknown executor: %s", spec.Executor)
		}

		d.Tasks[spec.ID] = Task{
			Id:       spec.ID,
			Name:     spec.Name,
			Executor: executor,
			Args:     spec.Args,
			Done:     make(chan bool), // Don't forget to initialize the Done channel
		}
	}

	// Fill the Edges map. The key is the ID of the task that depends on other tasks.
	// The value is a slice of Task objects that the task depends on.
	for taskId, taskIds := range taskDagSpec.Edges {
		for _, id := range taskIds {
			d.Edges[taskId] = append(d.Edges[taskId], d.Tasks[id])
		}
	}

	return nil
}

func (d *TaskDag) AddTask(task TaskSpec) {
	var executor Executor
	switch task.Executor {
	case "OpenAIAgent":
		executor = NewOpenAIAgent(task)
	case "GoogleTranslationAgent":
		executor = NewGoogleTranslationAgent(task)
	case "WebSocketClientAgent":
		executor = NewWebSocketClientAgent(task)
	case "OpenAIGenerativeAgent":
		executor = NewOpenAIGenerativeAgent(task)
	case "WebSocketServerAgent":
		executor = NewWebSocketServerAgent(task)
	case "DockerAgent":
		executor = NewDockerAgent(task)

	case "LocalFileSystemAgent":
		executor = NewLocalFileSystemAgent(task)
	default:
		return
	}

	d.Tasks[task.ID] = Task{
		Id:       task.ID,
		Name:     task.Name,
		Executor: executor,
		Args:     task.Args,
		Done:     make(chan bool),
	}
}

func (d *TaskDag) GetAllTasks() interface{} {
	return d.Tasks
}

func (d *TaskDag) GetTask(id string) (interface{}, interface{}) {
	return id, d.Tasks[id]
}

func (d *TaskDag) DeleteTask(id string) error {
	if _, ok := d.Tasks[id]; !ok {
		return fmt.Errorf("task %s not found", id)
	}
	delete(d.Tasks, id)
	return nil
}

func (d *TaskDag) UpdateTask(task Task) (Task, error) {
	if _, ok := d.Tasks[task.Id]; !ok {
		return Task{}, fmt.Errorf("task %s not found", task.Id)
	}
	t := d.Tasks[task.Id]
	t.UpdateWithMerge(task)
	return t, nil
}
