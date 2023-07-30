package main

import "context"

type TaskCrudRepository interface {
	CreateTask(ctx context.Context, task Task) (*Task, error)
	GetTasks(ctx context.Context) ([]Task, error)
	GetTask(ctx context.Context, id string) (*Task, error)
	DeleteTask(ctx context.Context, id string) error
	UpdateTask(ctx context.Context, id string, task Task) (*Task, error)
}
