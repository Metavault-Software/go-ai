package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"time"
)

type Workflow struct {
	ID          string
	Name        string
	Description string
	Tasks       map[string]Task
	CreatedAt   time.Time
	UpdatedAt   time.Time
	IsActive    bool
}

type WorkflowRepository interface {
	GetByID(ctx context.Context, userID string, workspaceID string, id string) (Workflow, error)
	Save(ctx context.Context, userID string, workspaceID string, workflow Workflow) error
	Delete(ctx context.Context, userID string, workspaceID string, id string) error
	GetAll(ctx context.Context, userID string, workspaceID string) ([]Workflow, error)
	Update(ctx context.Context, userID string, workspaceID string, workflow Workflow) error
}

func NewFirestoreWorkflowRepository(client *firestore.Client) WorkflowRepository {
	return &FirestoreWorkflowRepository{Client: client}
}

type FirestoreWorkflowRepository struct {
	Client *firestore.Client
}

func (r *FirestoreWorkflowRepository) GetByID(ctx context.Context, userID string, workspaceID string, id string) (Workflow, error) {
	userRef := r.Client.Collection("users").Doc(userID)
	workspaceRef := userRef.Collection("workspaces").Doc(workspaceID)
	doc, err := workspaceRef.Collection("workflows").Doc(id).Get(ctx)
	if err != nil {
		return Workflow{}, err
	}
	var workflow Workflow
	doc.DataTo(&workflow)
	return workflow, nil
}

func (r *FirestoreWorkflowRepository) Save(ctx context.Context, userID string, workspaceID string, workflow Workflow) error {
	userRef := r.Client.Collection("users").Doc(userID)
	workspaceRef := userRef.Collection("workspaces").Doc(workspaceID)
	_, _, err := workspaceRef.Collection("workflows").Add(ctx, workflow)
	return err
}

func (r *FirestoreWorkflowRepository) Delete(ctx context.Context, userID string, workspaceID string, id string) error {
	userRef := r.Client.Collection("users").Doc(userID)
	workspaceRef := userRef.Collection("workspaces").Doc(workspaceID)
	_, err := workspaceRef.Collection("workflows").Doc(id).Delete(ctx)
	return err
}

func (r *FirestoreWorkflowRepository) GetAll(ctx context.Context, userID string, workspaceID string) ([]Workflow, error) {
	userRef := r.Client.Collection("users").Doc(userID)
	workspaceRef := userRef.Collection("workspaces").Doc(workspaceID)
	iter := workspaceRef.Collection("workflows").Documents(ctx)
	var workflows []Workflow
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var workflow Workflow
		doc.DataTo(&workflow)
		workflows = append(workflows, workflow)
	}
	return workflows, nil
}

func (r *FirestoreWorkflowRepository) Update(ctx context.Context, userID string, workspaceID string, workflow Workflow) error {
	userRef := r.Client.Collection("users").Doc(userID)
	workspaceRef := userRef.Collection("workspaces").Doc(workspaceID)
	_, err := workspaceRef.Collection("workflows").Doc(workflow.ID).Set(ctx, workflow)
	return err
}
