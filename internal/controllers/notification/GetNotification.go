package notification

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
)

// GetNotifications retrieves all notifications for the authenticated user
func GetNotifications(c *gin.Context) {
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

    var notifications []models.Notification
    if err := database.DB.Where("user_id = ?", user.ID).Preload("Message").Preload("Sender").Find(&notifications).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}


