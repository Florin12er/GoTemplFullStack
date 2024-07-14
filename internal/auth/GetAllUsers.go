package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllUsers retrieves all users from the database
func GetAllUsers(c *gin.Context) {
	var users []models.User

	// Query the database for all users
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to retrieve users",
		})
		return
	}

	// Respond with the list of users
	c.JSON(http.StatusOK, users)
}

