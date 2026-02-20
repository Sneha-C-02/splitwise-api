package handlers

import (
	"net/http"
	"splitwise-api/config"
	"splitwise-api/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user with a bcrypt-hashed password
func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed. Email may already exist."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// GetUsers returns all registered users (password excluded)
func GetUsers(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)

	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"id":    u.ID,
			"name":  u.Name,
			"email": u.Email,
		})
	}

	c.JSON(http.StatusOK, gin.H{"users": result})
}