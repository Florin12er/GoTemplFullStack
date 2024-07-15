package notification

import(
    "github.com/gin-gonic/gin"
    "net/http"
    "strconv"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/Database"
)
// RemoveNotification deletes a specific notification
func RemoveNotification(c *gin.Context) {
    // Get the user ID from the JWT token
    userID, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Get the notification ID from the URL parameter
    notificationID, err := strconv.ParseUint(c.Param("notificationID"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
        return
    }

    // Find the notification
    var notification models.Notification
    if err := database.DB.First(&notification, notificationID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
        return
    }

    // Check if the notification belongs to the authenticated user
    if notification.UserID != uint(userID.(uint)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "You can only remove your own notifications"})
        return
    }

    // Delete the notification
    if err := database.DB.Delete(&notification).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete notification"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Notification removed successfully"})
}

