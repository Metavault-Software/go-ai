package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

func main() {
	r := gin.Default()

	// Add cors middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	homepage := r.Group("/")
	homepage.GET("/", WelcomeMessage)
	tasks := EmulateTasks()

	emulator := r.Group("/api/v1/emulator")
	{
		emulator.GET("/tasks", func(context *gin.Context) {
			context.JSON(200, tasks)
		})
	}

	store := AgentStore{store: make(map[string]Agent)}
	storeSpec := AgentStoreSpec{}
	err := storeSpec.LoadAgentsFromFile("agents.json", &store)
	if err != nil {
		log.Printf("Failed to load agents from file: %+v", err)
	}

	v1 := r.Group("/api/v1")
	{

		v1.POST("/agents", store.AddAgent)
		v1.GET("/agents", store.GetAgents)
		v1.GET("/agents/:agentId", store.GetAgent)
		v1.PUT("/agents/:agentId", store.UpdateAgent)
		v1.DELETE("/agents/:agentId", store.DeleteAgent)

		v1.POST("/tasks", CreateTask)
		v1.GET("/tasks", func(context *gin.Context) {
			context.JSON(200, tasks)
		})
		v1.GET("/tasks/:id", GetTask)
		v1.DELETE("/tasks/:id", DeleteTask)
		v1.PUT("/tasks/:id", UpdateTask)
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
	err = r.Run(":8080")
	if err != nil {
		log.Printf("Failed to start server: %+v", err)
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
