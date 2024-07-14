package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)
// In auth/handlers.go or a similar file in your auth package
func DeleteUserProfile(c *gin.Context) {
    // Get the user ID from the JWT token
    userID, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Fetch the existing user from the database
    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Delete the user from the database
    if err := database.DB.Delete(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User profile deleted successfully"})
}

