package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

type User struct {
	ID          string
	Email       string
	Password    string // make sure to hash the password before storing
	Workspaces  map[string]Workspace
	CreatedAt   time.Time
	UpdatedAt   time.Time
	FirstName   string
	LastName    string
	IsActive    bool
	LastLoginAt time.Time
}

type UserRepository interface {
	GetByID(ctx context.Context, id string) (User, error)
	Save(ctx context.Context, user User) (User, error)
	Delete(ctx context.Context, id string) error
	GetByEmail(ctx context.Context, email string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
	Update(ctx context.Context, user User) (User, error)
	GetClient() *firestore.Client
}

type FirestoreUserRepository struct {
	Client *firestore.Client
}

func (r *FirestoreUserRepository) GetClient() *firestore.Client {
	return r.Client
}

func NewFirestoreUserRepository() UserRepository {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "valid-actor-393616")
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	return &FirestoreUserRepository{Client: client}
}

func (r *FirestoreUserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	iter := r.Client.Collection("users").Where("Email", "==", email).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return User{}, err
	}
	var user User
	doc.DataTo(&user)
	return user, nil
}

func (r *FirestoreUserRepository) GetByID(ctx context.Context, id string) (User, error) {
	doc, err := r.Client.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return User{}, err
	}
	var user User
	doc.DataTo(&user)
	return user, nil
}

func (r *FirestoreUserRepository) Save(ctx context.Context, user User) (User, error) {
	_, _, err := r.Client.Collection("users").Add(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *FirestoreUserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.Client.Collection("users").Doc(id).Delete(ctx)
	return err
}

func (r *FirestoreUserRepository) GetAll(ctx context.Context) ([]User, error) {
	iter := r.Client.Collection("users").Documents(ctx)
	var users []User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user User
		doc.DataTo(&user)
		users = append(users, user)
	}
	return users, nil
}

func (r *FirestoreUserRepository) Update(ctx context.Context, user User) (User, error) {
	_, err := r.Client.Collection("users").Doc(user.ID).Set(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *FirestoreUserRepository) FindByEmail(ctx context.Context, email string) (User, error) {
	iter := r.Client.Collection("users").Where("Email", "==", email).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		return User{}, err
	}
	var user User
	doc.DataTo(&user)
	return user, nil
}
