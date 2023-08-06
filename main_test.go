package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestPassword(t *testing.T) {
	password := "secure_password"

	// Generate the bcrypt hash for the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating bcrypt hash:", err)
		return
	}

	fmt.Println("Bcrypt hash:", string(hashedPassword))
}

func TestTaskDag(t *testing.T) {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Open JSON file
	jsonFile, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Fatalf("Error closing JSON file: %v", err)
		}
	}(jsonFile)

	// Read JSON file
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Parse JSON into task specifications
	var taskGraph TaskGraph
	err = json.Unmarshal(byteValue, &taskGraph)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	started := make(chan Task)
	result := make(chan Task)
	go taskGraph.Run(started, result)

	for {
		select {
		case task, ok := <-started:
			if ok {
				fmt.Printf("Started: %s\n", task.Id)
			}
		case task, ok := <-result:
			if ok {
				fmt.Printf("Finished: %s\n", task.Id)
			} else {
				// No more tasks to process
				return
			}
		default:
			// Let's not hog the CPU
			time.Sleep(50 * time.Millisecond)
		}
	}
}

func TestTest1(t *testing.T) {
	ctx := context.Background()

	// Initialize Firestore client
	client, err := firestore.NewClient(ctx, "valid-actor-393616")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			t.Errorf("Failed to close client: %v", err)
		}
	}(client)

	const numUsers = 5
	const numWorkspaces = 5
	const numWorkflows = 5
	const numTasks = 5

	for i := 1; i <= numUsers; i++ {
		// Create a sample user
		user := User{
			ID:        fmt.Sprintf("User%d", i),
			Email:     fmt.Sprintf("user%d@example.com", i),
			Password:  "password", // Hash this in a real application
			FirstName: fmt.Sprintf("John%d", i),
			LastName:  fmt.Sprintf("Doe%d", i),
			IsActive:  true,
		}

		// Save the user
		_, err = client.Collection("users").Doc(user.ID).Set(ctx, user)
		if err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}

		for j := 1; j <= numWorkspaces; j++ {
			// Create a sample workspace
			workspace := Workspace{
				ID:          fmt.Sprintf("Workspace%d-%d", i, j),
				Name:        fmt.Sprintf("Workspace %d-%d", i, j),
				Description: fmt.Sprintf("Workspace %d-%d description", i, j),
				IsActive:    true,
				OwnerID:     user.ID,
			}

			// Save the workspace
			_, err = client.Collection("users").Doc(user.ID).Collection("workspaces").Doc(workspace.ID).Set(ctx, workspace)
			if err != nil {
				log.Fatalf("Failed to create workspace: %v", err)
			}

			for k := 1; k <= numWorkflows; k++ {
				// Create a sample workflow
				workflow := Workflow{
					ID:          fmt.Sprintf("Workflow%d-%d-%d", i, j, k),
					Name:        fmt.Sprintf("Workflow %d-%d-%d", i, j, k),
					Description: fmt.Sprintf("Workflow %d-%d-%d description", i, j, k),
					IsActive:    true,
				}

				// Save the workflow
				_, err = client.Collection("users").Doc(user.ID).Collection("workspaces").Doc(workspace.ID).Collection("workflows").Doc(workflow.ID).Set(ctx, workflow)
				if err != nil {
					log.Fatalf("Failed to create workflow: %v", err)
				}

				for l := 1; l <= numTasks; l++ {
					// Create a sample task
					task := Task{
						Id: fmt.Sprintf("Task%d-%d-%d-%d", i, j, k, l),
						Agent: Agent{
							AgentId: "1",
							Name:    fmt.Sprintf("Task %d-%d-%d-%d", i, j, k, l),
						},
						Labels: []string{"openai", "chat"},
						Args: map[string]interface{}{
							"messages": []map[string]interface{}{
								{
									"role":    "user",
									"content": fmt.Sprintf("Task %d-%d-%d-%d content", i, j, k, l),
								},
							},
						},
					}
					err := task.ExecutorFromExecutorType()
					if err != nil {
						t.Errorf("Failed to create task: %v", err)
					}

					// Save the task
					_, err = client.Collection("users").Doc(user.ID).Collection("workspaces").Doc(workspace.ID).Collection("workflows").Doc(workflow.ID).Collection("tasks").Doc(task.Id).Set(ctx, task)
					if err != nil {
						log.Fatalf("Failed to create task: %v", err)
					}
				}
			}
		}
	}

	log.Println("Sample data created successfully!")
}
