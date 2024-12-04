package handlers

import (
	"net/http"

	"github.com/valentinpelus/remediate/models"
	"github.com/valentinpelus/remediate/utils"

	"github.com/gin-gonic/gin"
)

// Function for logging in
func Login(c *gin.Context) {
	var user models.User

	// Check user credentials and generate a JWT token
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Check if credentials are valid (replace this logic with real authentication)
	if user.Username == "user" && user.Password == "password" {
		// Generate a JWT token
		token, err := utils.GenerateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}

// Function for registering a new user (for demonstration purposes)
func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Remember to securely hash passwords before storing them
	user.ID = 1 // Just for demonstration purposes
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}
