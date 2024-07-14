package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// EditUserProfile updates the profile information of the logged-in user
func EditUserProfile(c *gin.Context) {
	userID := utils.GetLoggedInUserID(
		c,
	) // Assume getLoggedInUserID is implemented as discussed earlier

	var update struct {
		UserName string `json:"username"`
		FullName string `json:"fullname"`
		Email    string `json:"email"`
	}
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Fetch the existing user record
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Update fields if they are provided
	if update.UserName != "" {
		user.UserName = update.UserName
	}
	if update.FullName != "" {
		user.FullName = update.FullName
	}
	if update.Email != "" {
		user.Email = update.Email
	}

	// Save the updated user record
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user profile updated successfully"})
}
