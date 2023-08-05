package main

import "context"

type TaskRepository interface {
	CreateTask(ctx context.Context, userID, workspaceID, workflowID string, task Task) (*Task, error)
	GetTasks(ctx context.Context, userID, workspaceID, workflowID string) ([]Task, error)
	GetByID(ctx context.Context, userID, workspaceID, workflowID, id string) (*Task, error)
	Delete(ctx context.Context, userID, workspaceID, workflowID, id string) error
	Update(ctx context.Context, userID, workspaceID, workflowID, id string, task Task) (*Task, error)
}
