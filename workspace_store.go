package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"time"
)

type Workspace struct {
	ID          string              `json:"id" firestore:"id"`
	Name        string              `json:"name" firestore:"name"`
	Description string              `json:"description" firestore:"description"`
	Workflows   map[string]Workflow `json:"workflows" firestore:"workflows"`
	CreatedAt   time.Time           `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" firestore:"updated_at"`
	IsActive    bool                `json:"is_active" firestore:"is_active"`
	OwnerID     string              `json:"owner_id" firestore:"owner_id"`
}

type WorkspaceRepository interface {
	GetByID(ctx context.Context, userID string, id string) (Workspace, error)
	Save(ctx context.Context, userID string, workspace Workspace) error
	Delete(ctx context.Context, userID string, id string) error
	GetAll(ctx context.Context, userID string) ([]Workspace, error)
	Update(ctx context.Context, userID string, workspace Workspace) error
}

type FirestoreWorkspaceRepository struct {
	Client *firestore.Client
}

func NewFirestoreWorkspaceRepository(client *firestore.Client) *FirestoreWorkspaceRepository {
	return &FirestoreWorkspaceRepository{Client: client}
}

func (r *FirestoreWorkspaceRepository) GetByID(ctx context.Context, userID string, id string) (Workspace, error) {
	userRef := r.Client.Collection("users").Doc(userID)
	doc, err := userRef.Collection("workspaces").Doc(id).Get(ctx)
	if err != nil {
		return Workspace{}, err
	}
	var workspace Workspace
	err = doc.DataTo(&workspace)
	if err != nil {
		return Workspace{}, err
	}
	return workspace, nil
}

func (r *FirestoreWorkspaceRepository) Save(ctx context.Context, userID string, workspace Workspace) error {
	userRef := r.Client.Collection("users").Doc(userID)
	_, _, err := userRef.Collection("workspaces").Add(ctx, workspace)
	return err
}

func (r *FirestoreWorkspaceRepository) Delete(ctx context.Context, userID string, id string) error {
	userRef := r.Client.Collection("users").Doc(userID)
	_, err := userRef.Collection("workspaces").Doc(id).Delete(ctx)
	return err
}

func (r *FirestoreWorkspaceRepository) GetAll(ctx context.Context, userID string) ([]Workspace, error) {
	userRef := r.Client.Collection("users").Doc(userID)
	iter := userRef.Collection("workspaces").Documents(ctx)
	var workspaces []Workspace
	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, err
		}
		var workspace Workspace
		err = doc.DataTo(&workspace)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, workspace)
	}
	return workspaces, nil
}

func (r *FirestoreWorkspaceRepository) Update(ctx context.Context, userID string, workspace Workspace) error {
	userRef := r.Client.Collection("users").Doc(userID)
	_, err := userRef.Collection("workspaces").Doc(workspace.ID).Set(ctx, workspace)
	return err
}
