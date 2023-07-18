package main

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGin(t *testing.T) {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.POST("/tasks", CreateTask)
		v1.GET("/tasks", GetAllTasks)
		v1.GET("/tasks/:id", GetTask)
		v1.DELETE("/tasks/:id", DeleteTask)
		v1.PUT("/tasks/:id", UpdateTask)
	}
	r.Run() // Listen and serve on 0.0.0.0:8080
}
