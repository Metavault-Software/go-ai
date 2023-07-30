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

func NewFirestoreTaskRepository() TaskCrudRepository {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "valid-actor-393616")
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	return &FirestoreTaskRepository{Client: client}
}

func (r *FirestoreTaskRepository) CreateTask(ctx context.Context, task Task) (*Task, error) {
	// Query Firestore to find the document with the given task ID
	iter := r.Client.Collection("workflow").Where("id", "==", task.Id).Documents(ctx)
	doc, err := iter.Next()

	// If a document with the given task ID is found, update it
	if !errors.Is(err, iterator.Done) {
		if err != nil {
			return nil, err
		}
		_, err = doc.Ref.Set(ctx, task)
		if err != nil {
			return nil, err
		}
	} else {
		// If no document with the given task ID is found, create a new one
		_, _, err := r.Client.Collection("workflow").Add(ctx, task)
		if err != nil {
			return nil, err
		}
	}

	return &task, nil
}

func (r *FirestoreTaskRepository) GetTasks(ctx context.Context) ([]Task, error) {
	iter := r.Client.Collection("workflow").Documents(ctx)
	var tasks []Task
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var task Task
		doc.DataTo(&task)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *FirestoreTaskRepository) GetTask(ctx context.Context, id string) (*Task, error) {
	doc, err := r.Client.Collection("workflow").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var task Task
	doc.DataTo(&task)
	return &task, nil
}

func (r *FirestoreTaskRepository) DeleteTask(ctx context.Context, id string) error {
	// Query Firestore to find the document with the given task ID
	iter := r.Client.Collection("workflow").Where("id", "==", id).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return err
	}
	// Delete the document
	_, err = doc.Ref.Delete(ctx)
	return err
}

func (r *FirestoreTaskRepository) UpdateTask(ctx context.Context, id string, task Task) (*Task, error) {
	// Query Firestore to find the document with the given task ID
	iter := r.Client.Collection("workflow").Where("id", "==", id).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return nil, err
	}

	// Update the document
	_, err = doc.Ref.Set(ctx, task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
