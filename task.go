package main

import (
	"context"
	"fmt"
	"time"
)

type Task struct {
	Agent
	Id           string                 `json:"id" firestore:"id"`
	Args         map[string]interface{} `json:"args" firestore:"args"`
	Executor     Executor               `json:"-" firestore:"-"`
	Status       chan TaskStatus        `json:"-" firestore:"-"`
	Result       interface{}            `json:"result" firestore:"result"`
	Error        error                  `json:"error" firestore:"-"`
	Labels       []string               `json:"labels" firestore:"labels"`
	StartTime    time.Time              `json:"start_time" firestore:"start_time"`
	EndTime      time.Time              `json:"end_time" firestore:"end_time"`
	EstDuration  time.Duration          `json:"estimated_duration" firestore:"estimated_duration"`
	CompDuration time.Duration          `json:"completed_duration" firestore:"completed_duration"`
	Done         chan bool              `json:"-" firestore:"-"`
}

type TaskStatus int

const (
	Pending TaskStatus = iota
	Running
	Completed
	Failed
	Cancelled
)

type Executor interface {
	Execute(ctx context.Context, task *Task) error
}

func (t *Task) Run() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		// Update task status to Running
		t.Status <- Running
		// Note the start time
		t.StartTime = time.Now()
		// Execute the task
		err := t.Executor.Execute(ctx, t)
		// Note the end time
		t.EndTime = time.Now()
		// Calculate the completed duration
		t.CompDuration = t.EndTime.Sub(t.StartTime)
		// Check for errors from Execute
		if err != nil {
			if err == context.Canceled {
				t.Status <- Cancelled
				fmt.Println("Job was cancelled")
			} else {
				t.Status <- Failed
				t.Error = err
				fmt.Println("Job failed with error:", err)
			}
		} else {
			// If there was no error, mark the task as Completed
			t.Status <- Completed
			fmt.Println("Job completed successfully")
		}
	}()

	return cancel
}

func (t *Task) UpdateWithMerge(newTask Task) *Task {
	if newTask.Id != "" {
		t.Id = newTask.Id
	}
	if newTask.Name != "" {
		t.Name = newTask.Name
	}
	if newTask.Result != nil {
		t.Result = newTask.Result
	}
	if newTask.Error != nil {
		t.Error = newTask.Error
	}
	if newTask.Executor != nil {
		t.Executor = newTask.Executor
	}
	if newTask.Labels != nil {
		if t.Labels == nil {
			t.Labels = make([]string, 0)
		}
		for key, value := range newTask.Labels {
			t.Labels[key] = value
		}
	}
	if !newTask.StartTime.IsZero() {
		t.StartTime = newTask.StartTime
	}
	if !newTask.EndTime.IsZero() {
		t.EndTime = newTask.EndTime
	}
	if newTask.EstDuration != 0 {
		t.EstDuration = newTask.EstDuration
	}
	if newTask.CompDuration != 0 {
		t.CompDuration = newTask.CompDuration
	}
	if newTask.Args != nil {
		if t.Args == nil {
			t.Args = make(map[string]interface{})
		}
		for key, value := range newTask.Args {
			t.Args[key] = value
		}
	}
	if newTask.Done != nil {
		t.Done = newTask.Done
	}
	return t
}

func (t *Task) ExecutorFromExecutorType() error {
	switch t.ExecutorType {
	case "OpenAIAgent":
		t.Executor = NewOpenAIAgent(*t)
	case "GoogleTranslationAgent":
		t.Executor = NewGoogleTranslationAgent(*t)
	case "WebSocketClientAgent":
		t.Executor = NewWebSocketClientAgent(*t)
	case "OpenAIGenerativeAgent":
		t.Executor = NewOpenAIGenerativeAgent(*t)
	case "WebSocketServerAgent":
		t.Executor = NewWebSocketServerAgent(*t)
	case "DockerAgent":
		t.Executor = NewDockerAgent(*t)
	case "LocalFileSystemAgent":
		t.Executor = NewLocalFileSystemAgent(*t)
	default:
		return fmt.Errorf("unknown executor: %s", t.Executor)
	}
	return nil
}
