package message

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// DeleteMessage allows a user to delete a message
func DeleteMessage(c *gin.Context) {
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

	// Get the message ID from the URL parameter
	messageID, err := strconv.ParseUint(c.Param("messageID"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	// Fetch the message from the database
	var message models.Message
	if err := database.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// Check if the user is the sender of the message
	if message.SenderID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own messages"})
		return
	}

	// Delete the message
	if err := database.DB.Delete(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message deleted successfully"})
}
