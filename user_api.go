package main

import (
	"context"
	"github.com/gin-gonic/gin"
)

type UserWebApi struct {
	Repo UserRepository
}

func NewUserApi(repo UserRepository) *UserWebApi {
	return &UserWebApi{Repo: NewFirestoreUserRepository()}
}

func (ua *UserWebApi) CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	newUser, err := ua.Repo.Save(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add user"})
		return
	}
	c.JSON(201, newUser)
}

func (ua *UserWebApi) GetUsers(c *gin.Context) {
	users, err := ua.Repo.GetAll(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get users"})
		return
	}
	c.JSON(200, users)
}

func (ua *UserWebApi) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := ua.Repo.GetByID(context.Background(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, user)
}

func (ua *UserWebApi) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := ua.Repo.Delete(context.Background(), id)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted successfully"})
}

func (ua *UserWebApi) UpdateUser(c *gin.Context) {
	//id := c.Param("id")
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	updatedUser, err := ua.Repo.Update(context.Background(), user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}
	c.JSON(200, updatedUser)
}

type UserApiInterface interface {
	CreateUser(c *gin.Context)
	GetUsers(c *gin.Context)
	GetUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}
