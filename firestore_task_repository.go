package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"log"
)

type FirestoreTaskRepository struct {
	Client *firestore.Client
}

func (r *FirestoreTaskRepository) GetByID(
	ctx context.Context,
	userID, workspaceID, workflowID, id string,
) (*Task, error) {
	doc, err := r.Client.Collection("users").Doc(userID).
		Collection("workspaces").Doc(workspaceID).
		Collection("workflows").Doc(workflowID).
		Collection("tasks").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var task Task
	err = doc.DataTo(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func NewFirestoreTaskRepository() TaskRepository {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "valid-actor-393616")
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	return &FirestoreTaskRepository{Client: client}
}

func (r *FirestoreTaskRepository) CreateTask(
	ctx context.Context,
	userID string,
	workspaceID string,
	workflowID string,
	task Task,
) (*Task, error) {
	// Path to the specific user's workspace's workflow's tasks collection
	path := r.Client.Collection("users").Doc(userID).
		Collection("workspaces").Doc(workspaceID).
		Collection("workflows").Doc(workflowID).
		Collection("tasks")

	// Add a new task to the tasks collection
	docRef, _, err := path.Add(ctx, task)
	if err != nil {
		return nil, err
	}

	task.Id = docRef.ID
	return &task, nil
}

func (r *FirestoreTaskRepository) GetTasks(
	ctx context.Context,
	userID string,
	workspaceID string,
	workflowID string,
) ([]Task, error) {
	iter := r.Client.Collection("users").Doc(userID).
		Collection("workspaces").Doc(workspaceID).
		Collection("workflows").Doc(workflowID).
		Collection("tasks").Documents(ctx)
	var tasks []Task
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		var task Task
		err = doc.DataTo(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *FirestoreTaskRepository) Update(
	ctx context.Context,
	userID string,
	workspaceID string,
	workflowID string,
	id string,
	task Task,
) (*Task, error) {
	_, err := r.Client.Collection("users").Doc(userID).
		Collection("workspaces").Doc(workspaceID).
		Collection("workflows").Doc(workflowID).
		Collection("tasks").Doc(id).Set(ctx, task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *FirestoreTaskRepository) Delete(
	ctx context.Context,
	userID string,
	workspaceID string,
	workflowID string,
	id string,
) error {
	_, err := r.Client.Collection("users").Doc(userID).
		Collection("workspaces").Doc(workspaceID).
		Collection("workflows").Doc(workflowID).
		Collection("tasks").Doc(id).Delete(ctx)
	return err
}
