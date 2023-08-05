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

//func (ua *UserWebApi) Login(c *gin.Context) {
//	// ... Same code as before
//
//	// Create a session token
//	sessionToken := generateSessionToken()
//
//	// Store the session token in Firestore
//	_, err := ua.Repo.GetClient().Collection("sessions").Doc(sessionToken).Set(context.Background(), map[string]interface{}{
//		"userID":    user.ID,
//		"createdAt": time.Now(),
//		"expiresAt": time.Now().Add(24 * time.Hour), // Example expiration time
//	})
//	if err != nil {
//		c.JSON(500, gin.H{"error": "Failed to create session"})
//		return
//	}
//
//	c.JSON(200, gin.H{"message": "Logged in successfully", "token": sessionToken})
//}
//
//func generateSessionToken() string {
//	// Implement your logic to generate a unique session token
//}
//
//func (ua *UserWebApi) Logout(c *gin.Context) {
//	// Retrieve the session token from request, e.g., from a header
//	sessionToken := c.GetHeader("Authorization")
//
//	// Delete the session token from Firestore
//	_, err := ua.FirestoreClient.Collection("sessions").Doc(sessionToken).Delete(context.Background())
//	if err != nil {
//		c.JSON(500, gin.H{"error": "Failed to log out"})
//		return
//	}
//
//	c.JSON(200, gin.H{"message": "Logged out successfully"})
//}
