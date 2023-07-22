package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", CreateTask)
		//v1.GET("/tasks", GetAllTasks)
		v1.GET("/tasks", EmulateAllTasks)
		v1.GET("/tasks/:id", GetTask)
		v1.DELETE("/tasks/:id", DeleteTask)
		v1.PUT("/tasks/:id", UpdateTask)
	}
	r.Run(":8080")
}

func WelcomeMessage(context *gin.Context) {
	context.JSON(200, gin.H{"message": "Welcome to the Task API"})
}

func EmulateAllTasks(context *gin.Context) {
	emulator := TaskDagEmulator{
		taskDag,
	}
	tasks := emulator.Generate100RandomTasks()
	context.JSON(200, tasks)
}
