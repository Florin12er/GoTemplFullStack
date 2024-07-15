package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)


func ResetRequest(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email exists in the database
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	// Generate and store reset code
	code := utils.GenerateResetCode()
	utils.ResetCodes[req.Email] = utils.ResetCode{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// Send email with reset code
	if err := utils.SendResetEmail(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to email"})
}

func ResetPassword(c *gin.Context) {
    var req struct {
        Email       string `json:"email" binding:"required,email"`
        Code        string `json:"code" binding:"required,min=6"`
        NewPassword string `json:"newPassword" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Verify reset code
    resetCode, exists := utils.ResetCodes[req.Email]
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No reset code found for this email"})
        return
    }
    if resetCode.Code != req.Code {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset code"})
        return
    }
    if time.Now().After(resetCode.ExpiresAt) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Reset code has expired"})
        return
    }

    // Update password in the database
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    if err := database.DB.Model(&models.User{}).Where("email = ?", req.Email).Update("password", string(hashedPassword)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    // Remove used reset code
    delete(utils.ResetCodes, req.Email)

    c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

