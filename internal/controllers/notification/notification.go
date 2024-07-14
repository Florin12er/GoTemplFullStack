package notification

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetNotifications retrieves all notifications for a user
func GetNotifications(c *gin.Context) {
	userID := utils.GetLoggedInUserID(c) // Define how to get the logged-in user ID
	var notifications []models.Notification
	if err := database.DB.Preload("Message").Preload("Sender").
		Where("user_id = ?", userID).
		Order("created_at DESC").Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve notifications"})
		return
	}

	// Mark notifications as read
	for _, notification := range notifications {
		if !notification.IsRead {
			notification.IsRead = true
			database.DB.Save(&notification)
		}
	}

	c.JSON(http.StatusOK, notifications)
}

// RemoveNotification removes a notification for a user
func RemoveNotification(c *gin.Context) {
	notificationID, err := strconv.Atoi(c.Param("notificationID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification ID"})
		return
	}

	// Check if the notification exists
	var notification models.Notification
	if err := database.DB.First(&notification, notificationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	// Delete the notification
	if err := database.DB.Delete(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification removed successfully"})
}
