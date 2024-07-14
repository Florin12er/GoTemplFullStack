package auth

import (
	"net/http"

	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type EditUserRequest struct {
	FullName    string `json:"fullName"`
	Email       string `json:"email"`
	UserName    string `json:"userName"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func EditUserProfile(c *gin.Context) {
	// Get the user from the context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, ok := userInterface.(models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
		return
	}

	// Parse the request body
	var req EditUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.UserName != "" {
		user.UserName = req.UserName
	}

	// Handle password change
	if req.OldPassword != "" && req.NewPassword != "" {
		// Verify old password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}

		// Hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(req.NewPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		user.Password = string(hashedPassword)
	}

	// Save the updated user to the database
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Return the updated user (excluding the password)
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{"message": "User profile updated successfully", "user": user})
}
