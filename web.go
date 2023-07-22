package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var taskDag = TaskDag{
	Tasks: make(map[string]Task),
	Edges: make(map[string][]Task),
}

// CreateTask Handler for POST /tasks
func CreateTask(c *gin.Context) {
	var task TaskSpec
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// Add the task to the TaskDag
	taskDag.AddTask(task)
	c.JSON(201, task)
}

// GetAllTasks Handler for GET /tasks
func GetAllTasks(c *gin.Context) {
	tasks := taskDag.GetAllTasks()
	fmt.Printf("%+v\n", tasks)
	c.JSON(200, tasks)
}

// GetTask Handler for GET /tasks/:id
func GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := taskDag.GetTask(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, task)
}

// Handler for DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := taskDag.DeleteTask(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}

// Handler for PUT /tasks/:id
func UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	task.Id = id // Ensure the task ID matches the one in the path
	updateTask, err2 := taskDag.UpdateTask(task)
	if err2 != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, updateTask)
}

func EmulateAllTasks() []Task {
	emulator := TaskDagEmulator{
		taskDag,
	}
	tasks := emulator.Generate100RandomTasks()
	return tasks
}
