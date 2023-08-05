package main

import (
	"context"
	"github.com/gin-gonic/gin"
)

type WorkflowApi struct {
	Repo WorkflowRepository
}

func NewWorkflowApi(repo WorkflowRepository) *WorkflowApi {
	return &WorkflowApi{Repo: repo}
}

func (wa *WorkflowApi) GetWorkflows(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflows, err := wa.Repo.GetAll(context.Background(), userID, workspaceID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get workflows"})
		return
	}
	c.JSON(200, workflows)
}

func (wa *WorkflowApi) GetWorkflow(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("id")
	workflow, err := wa.Repo.GetByID(context.Background(), userID, workspaceID, workflowID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Workflow not found"})
		return
	}
	c.JSON(200, workflow)
}

func (wa *WorkflowApi) CreateWorkflow(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	var workflow Workflow
	if err := c.ShouldBindJSON(&workflow); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := wa.Repo.Save(context.Background(), userID, workspaceID, workflow)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add workflow"})
		return
	}
	c.JSON(201, workflow)
}

func (wa *WorkflowApi) DeleteWorkflow(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	workflowID := c.Param("id")
	err := wa.Repo.Delete(context.Background(), userID, workspaceID, workflowID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Workflow not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Workflow deleted successfully"})
}

func (wa *WorkflowApi) UpdateWorkflow(c *gin.Context) {
	userID := c.Param("user_id")
	workspaceID := c.Param("workspace_id")
	var workflow Workflow
	if err := c.ShouldBindJSON(&workflow); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := wa.Repo.Update(context.Background(), userID, workspaceID, workflow)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update workflow"})
		return
	}
	c.JSON(200, workflow)
}
