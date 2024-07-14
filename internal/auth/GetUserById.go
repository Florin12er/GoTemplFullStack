package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// GetUserByID retrieves a user by their ID
func GetUserByID(c *gin.Context) {
	idStr := c.Param("userID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

