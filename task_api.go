package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

type TaskApi struct {
	Repo TaskCrudRepository
}

func NewTaskApi(repo TaskCrudRepository) *TaskApi {
	return &TaskApi{Repo: repo}
}

func (ts *TaskApi) CreateTask(c *gin.Context) {
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := task.ExecutorFromExecutorType()
	if err != nil {
		log.Printf("Failed to create task: %+v", err)
	}
	newTask, err := ts.Repo.CreateTask(context.Background(), task)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add task"})
		return
	}
	c.JSON(201, newTask)
}

func (ts *TaskApi) GetTasks(c *gin.Context) {
	tasks, err := ts.Repo.GetTasks(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get tasks"})
		return
	}
	c.JSON(200, tasks)
}

func (ts *TaskApi) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := ts.Repo.GetTask(context.Background(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(200, task)
}

func (ts *TaskApi) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := ts.Repo.DeleteTask(context.Background(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}

func (ts *TaskApi) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	updatedTask, err := ts.Repo.UpdateTask(context.Background(), id, task)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update task"})
		return
	}
	c.JSON(200, updatedTask)
}

//func (ts *TaskApi) EmulateTasks() []Task {
//	emulator := TaskDagEmulator{
//		taskDag,
//	}
//	tasks := emulator.Generate100RandomTasks()
//	return tasks
//}
