package message

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/Database"
    "strconv"
)

// EditMessage allows a user to edit a message
func EditMessage(c *gin.Context) {
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

    // Parse the request body
    var req struct {
        Text string `json:"text" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
        c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own messages"})
        return
    }

    // Update the message
    message.Text = req.Text
    if err := database.DB.Save(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully", "data": message})
}
