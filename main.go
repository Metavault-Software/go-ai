package main

import (
	"encoding/json"
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
	tasks := EmulateAllTasks()

	emulator := r.Group("/api/v1/emulator")
	{
		emulator.GET("/tasks", func(context *gin.Context) {
			context.JSON(200, tasks)
		})
	}

	v1 := r.Group("/api/v1")
	{
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
						v := struct {
							Status int `json:"status"`
						}{
							Status: int(Completed),
						}
						err = conn.WriteJSON(v)
						if err != nil {
							fmt.Printf("Failed to write JSON: %+v", err)
							break
						}
						err = conn.Close()
						if err != nil {
							fmt.Printf("Failed to close websocket connection: %+v", err)
						}
						return
					case status := <-task.Status:
						v := struct {
							Status int `json:"status"`
						}{
							Status: int(status),
						}
						err = conn.WriteJSON(v)
						if err != nil {
							fmt.Printf("Failed to write JSON: %+v", err)
							break
						}
					}
				}
			}
			v1.GET(statusPath, handlers)
		}
	}
	err := r.Run(":8080")
	if err != nil {
		log.Printf("Failed to start server: %+v", err)
	}
}

func ToJson(status TaskStatus) ([]byte, error) {
	marshal, err := json.Marshal(struct {
		Status int `json:"status"`
	}{
		Status: int(status),
	})
	return marshal, err
}

func WelcomeMessage(context *gin.Context) {
	context.JSON(200, gin.H{"message": "Welcome to the Task API"})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
