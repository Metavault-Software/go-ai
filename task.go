package main

import (
	"context"
	"fmt"
	"time"
)

type Task struct {
	Id           string                 // Unique identifier for each task
	Name         string                 // Job name for better understanding
	Status       chan TaskStatus        `json:"-"` // Status of the task: Pending, Running, Completed, Failed, etc.
	Result       interface{}            // Result data after task execution
	Error        error                  // Any error encountered during task execution
	Executor     Executor               // An interface that knows how to execute the task
	Labels       map[string]interface{} // Additional task details as key-value pairs
	StartTime    time.Time              // Time when the task started
	EndTime      time.Time              // Time when the task ended
	EstDuration  time.Duration          // Estimated duration of the task
	CompDuration time.Duration          // Completed duration of the task
	Args         map[string]interface{} // Arguments to be passed to the executor
	Done         chan bool              `json:"-"` // Channel to signal task completion
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

func (d *Task) UpdateWithMerge(newTask Task) *Task {
	if newTask.Id != "" {
		d.Id = newTask.Id
	}
	if newTask.Name != "" {
		d.Name = newTask.Name
	}
	if newTask.Result != nil {
		d.Result = newTask.Result
	}
	if newTask.Error != nil {
		d.Error = newTask.Error
	}
	if newTask.Executor != nil {
		d.Executor = newTask.Executor
	}
	if newTask.Labels != nil {
		if d.Labels == nil {
			d.Labels = make(map[string]interface{})
		}
		for key, value := range newTask.Labels {
			d.Labels[key] = value
		}
	}
	if !newTask.StartTime.IsZero() {
		d.StartTime = newTask.StartTime
	}
	if !newTask.EndTime.IsZero() {
		d.EndTime = newTask.EndTime
	}
	if newTask.EstDuration != 0 {
		d.EstDuration = newTask.EstDuration
	}
	if newTask.CompDuration != 0 {
		d.CompDuration = newTask.CompDuration
	}
	if newTask.Args != nil {
		if d.Args == nil {
			d.Args = make(map[string]interface{})
		}
		for key, value := range newTask.Args {
			d.Args[key] = value
		}
	}
	if newTask.Done != nil {
		d.Done = newTask.Done
	}
	return d
}
