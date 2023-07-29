package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var taskDag = TaskGraph{
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
	// Add the task to the TaskGraph
	taskDag.AddTask(task)
	c.JSON(201, task)
}

// GetTasks Handler for GET /tasks
func GetTasks(c *gin.Context) {
	tasks := taskDag.GetTasks()
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

// DeleteTask Handler for DELETE /tasks/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := taskDag.DeleteTask(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}

// UpdateTask Handler for PUT /tasks/:id
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

func EmulateTasks() []Task {
	emulator := TaskDagEmulator{
		taskDag,
	}
	tasks := emulator.Generate100RandomTasks()
	return tasks
}
