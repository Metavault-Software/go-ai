package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	r := gin.Default()

	// Add cors middleware
	SetupCors(r)

	SetupHomeScreen(r)

	taskApi, agentApi, err := SetupDatastore()

	SetupRoutes(r, agentApi, taskApi)

	err = r.Run(":8080")
	if err != nil {
		log.Printf("Failed to start server: %+v", err)
	}
}

func SetupDatastore() (*TaskApi, AgentStore, error) {
	repository := NewFirestoreTaskRepository()
	taskApi := NewTaskApi(repository)
	store := AgentStore{store: make(map[string]Agent)}
	storeSpec := AgentStoreSpec{}
	err := storeSpec.LoadAgentsFromFile("agents.json", &store)
	if err != nil {
		log.Printf("Failed to load agents from file: %+v", err)
	}
	return taskApi, store, err
}

func SetupHomeScreen(r *gin.Engine) {
	homepage := r.Group("/")
	homepage.GET("/", WelcomeMessage)
}

func SetupRoutes(r *gin.Engine, store AgentStore, taskApi *TaskApi) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "valid-actor-393616")
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	v1 := r.Group("/api/v1")
	{
		SetupTaskRoutes(v1, taskApi)
		SetupWorkflowRoutes(v1, NewWorkflowApi(NewFirestoreWorkflowRepository(client)))
		SetupWorkspaceRoutes(v1, NewWorkspaceApi(NewFirestoreWorkspaceRepository(client)))
		SetupUserRoutes(v1, NewUserApi(NewFirestoreUserRepository()))
		SetupAgentRoutes(v1, store)
		var tasks []Task
		SetupStatus(tasks, v1)
	}
}

func SetupWorkspaceRoutes(v1 *gin.RouterGroup, workspaceApi *WorkspaceApi) {
	v1.POST("/users/:user_id/workspaces", workspaceApi.CreateWorkspace)
	v1.GET("/users/:user_id/workspaces", workspaceApi.GetWorkspaces)
	v1.GET("/users/:user_id/workspaces/:workspace_id", workspaceApi.GetWorkspace)
	v1.PUT("/users/:user_id/workspaces/:workspace_id", workspaceApi.UpdateWorkspace)
	v1.DELETE("/users/:user_id/workspaces/:workspace_id", workspaceApi.DeleteWorkspace)
}

func SetupWorkflowRoutes(v1 *gin.RouterGroup, workflowApi *WorkflowApi) {
	v1.POST("/users/:user_id/workspaces/:workspace_id/workflows", workflowApi.CreateWorkflow)
	v1.GET("/users/:user_id/workspaces/:workspace_id/workflows", workflowApi.GetWorkflows)
	v1.GET("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id", workflowApi.GetWorkflow)
	v1.PUT("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id", workflowApi.UpdateWorkflow)
	v1.DELETE("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id", workflowApi.DeleteWorkflow)
}

func SetupTaskRoutes(v1 *gin.RouterGroup, taskApi *TaskApi) {
	v1.POST("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id/tasks", taskApi.CreateTask)
	v1.GET("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id/tasks", taskApi.GetTasks)
	v1.GET("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id/tasks/:id", taskApi.GetTask)
	v1.PUT("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id/tasks/:id", taskApi.UpdateTask)
	v1.DELETE("/users/:user_id/workspaces/:workspace_id/workflows/:workflow_id/tasks/:id", taskApi.DeleteTask)
}

func SetupUserRoutes(v1 *gin.RouterGroup, userApi *UserWebApi) {
	v1.POST("/users", userApi.CreateUser)
	v1.GET("/users", userApi.GetUsers)
	v1.GET("/users/:user_id", userApi.GetUser)
	v1.PUT("/users/:user_id", userApi.UpdateUser)
	v1.DELETE("/users/:user_id", userApi.DeleteUser)
}

func SetupAgentRoutes(v1 *gin.RouterGroup, store AgentStore) {
	v1.POST("/agents", store.AddAgent)
	v1.GET("/agents", store.GetAgents)
	v1.GET("/agents/:agentId", store.GetAgent)
	v1.PUT("/agents/:agentId", store.UpdateAgent)
	v1.DELETE("/agents/:agentId", store.DeleteAgent)
}

func SetupCors(r *gin.Engine) gin.IRoutes {
	return r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))
}

func SetupStatus(tasks []Task, v1 *gin.RouterGroup) {
	for _, task := range tasks {
		task := task // Create a new 'task' variable in this scope, otherwise all goroutines will share the same loop variable
		statusPath := "/ws/" + task.Id + "/status"
		handlers := func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				fmt.Printf("Failed to set websocket upgrade: %+v", err)
				return
			}

			defer func(conn *websocket.Conn) {
				err := conn.Close()
				if err != nil {
					log.Printf("Failed to close websocket connection: %+v", err)
				}
			}(conn)

			for {
				select {
				case <-task.Done:
					WriteDone(err, conn)
					return
				case status := <-task.Status:
					WriteStatus(status, err, conn)
				}
			}
		}
		v1.GET(statusPath, handlers)
	}
}

func WriteDone(err error, conn *websocket.Conn) {
	type TaskStatus struct {
		Status int `json:"status"`
	}
	err = conn.WriteJSON(TaskStatus{
		Status: int(Completed),
	})
	if err != nil {
		fmt.Printf("Failed to write JSON: %+v", err)
	}
	err = conn.Close()
	if err != nil {
		fmt.Printf("Failed to close websocket connection: %+v", err)
	}
	return
}

func WriteStatus(status TaskStatus, err error, conn *websocket.Conn) {
	err = conn.WriteJSON(struct {
		Status int `json:"status"`
	}{
		Status: int(status),
	})
	if err != nil {
		fmt.Printf("Failed to write JSON: %+v", err)
	}
}

func WelcomeMessage(context *gin.Context) {
	context.JSON(200, gin.H{"message": "Welcome to the Task API"})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
