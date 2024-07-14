package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResetPassword resets the user's password
func ResetPassword(c *gin.Context) {
	var body struct {
		FullName    string `json:"full_name"`
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	// Bind JSON request body into the defined struct
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Validate that all required fields are provided
	if body.FullName == "" || body.Email == "" || body.NewPassword == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "full_name, email, and new_password are required fields"},
		)
		return
	}

	// Retrieve the user from the database based on full name and email
	var user models.User
	if err := database.DB.Where("full_name = ? AND email = ?", body.FullName, body.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Update the user's password
	user.Password = body.NewPassword
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reset password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password reset successfully"})
}
