package main

import (
	"context"
	base64 "encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
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

	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	fmt.Println("Hashed Password:", string(hashedPassword)) // Add this line to print the hashed password

	user.Password = string(hashedPassword)

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

// generateSessionToken generates a unique session token.
func generateSessionToken() string {
	// Define the length of the session token in bytes
	tokenLength := 32

	// Generate random bytes
	randomBytes := make([]byte, tokenLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		// Handle error, if any
		panic(err)
	}

	// Encode the random bytes to base64 string
	token := base64.RawStdEncoding.EncodeToString(randomBytes)

	return token
}

func (ua *UserWebApi) Login(c *gin.Context) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := ua.Repo.GetByEmail(context.Background(), credentials.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Validate password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a session token
	sessionToken := generateSessionToken()
	user.SessionToken = sessionToken

	// Save the updated user with the session token
	_, err = ua.Repo.Update(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"session_token": sessionToken})
}

func (ua *UserWebApi) Logout(c *gin.Context) {
	sessionToken := c.GetHeader("Authorization")
	if sessionToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session token not provided"})
		return
	}

	// Get the user with the given session token
	user, err := ua.Repo.GetBySessionToken(context.Background(), sessionToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session token"})
		return
	}

	// Clear the session token
	user.SessionToken = ""

	// Save the updated user without the session token
	_, err = ua.Repo.Update(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
