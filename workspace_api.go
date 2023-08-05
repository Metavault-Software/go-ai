package main

import (
	"context"
	"github.com/gin-gonic/gin"
)

type WorkspaceApi struct {
	Repo WorkspaceRepository
}

func NewWorkspaceApi(repo WorkspaceRepository) *WorkspaceApi {
	return &WorkspaceApi{Repo: repo}
}

func (wsa *WorkspaceApi) GetWorkspaces(c *gin.Context) {
	param := c.Param("user_id")
	workspaces, err := wsa.Repo.GetAll(context.Background(), param)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get workspaces"})
		return
	}
	c.JSON(200, workspaces)
}

func (wsa *WorkspaceApi) GetWorkspace(c *gin.Context) {
	id := c.Param("id")
	userId := c.Param("user_id")
	workspace, err := wsa.Repo.GetByID(context.Background(), userId, id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Workspace not found"})
		return
	}
	c.JSON(200, workspace)
}

func (wsa *WorkspaceApi) CreateWorkspace(c *gin.Context) {
	userID := c.Param("user_id") // retrieve the user ID from the request context
	var workspace Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := wsa.Repo.Save(context.Background(), userID, workspace) // pass user ID to repository
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add workspace"})
		return
	}
	c.JSON(201, workspace)
}

func (wsa *WorkspaceApi) DeleteWorkspace(c *gin.Context) {
	userID := c.Param("user_id")
	id := c.Param("id")
	err := wsa.Repo.Delete(context.Background(), userID, id) // pass user ID to repository
	if err != nil {
		c.JSON(404, gin.H{"error": "Workspace not found"})
		return
	}
	c.JSON(200, gin.H{"message": "Workspace deleted successfully"})
}

func (wsa *WorkspaceApi) UpdateWorkspace(c *gin.Context) {
	userID := c.Param("user_id")
	var workspace Workspace
	if err := c.ShouldBindJSON(&workspace); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	err := wsa.Repo.Update(context.Background(), userID, workspace) // pass user ID to repository
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update workspace"})
		return
	}
	c.JSON(200, workspace)
}
