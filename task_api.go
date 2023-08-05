package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
)

type TaskApi struct {
	Repo TaskRepository
}

func NewTaskApi(repo TaskRepository) *TaskApi {
	return &TaskApi{Repo: repo}
}

// ...

func (ts *TaskApi) CreateTask(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("workflow_id")

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := task.ExecutorFromExecutorType()
	if err != nil {
		log.Printf("Failed to create task: %+v", err)
	}

	newTask, err := ts.Repo.CreateTask(context.Background(), userID, workspaceID, workflowID, task)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add task"})
		return
	}
	c.JSON(201, newTask)
}

func (ts *TaskApi) GetTasks(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("workflow_id")

	tasks, err := ts.Repo.GetTasks(context.Background(), userID, workspaceID, workflowID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get tasks"})
		return
	}
	c.JSON(200, tasks)
}

func (ts *TaskApi) GetTask(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("workflow_id")
	id := c.Param("id")

	task, err := ts.Repo.GetByID(context.Background(), userID, workspaceID, workflowID, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(200, task)
}

func (ts *TaskApi) DeleteTask(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("workflow_id")
	id := c.Param("id")

	err := ts.Repo.Delete(context.Background(), userID, workspaceID, workflowID, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Task deleted successfully"})
}

func (ts *TaskApi) UpdateTask(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("workflow_id")
	id := c.Param("id")

	var task Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	updatedTask, err := ts.Repo.Update(context.Background(), userID, workspaceID, workflowID, id, task)
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
