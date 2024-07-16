package message

import (
    "GoMessageApp/internal/Database"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/templates"
    "context"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

func SendMessage(c *gin.Context) {
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

    // Parse the form data
    receiverID, err := strconv.ParseUint(c.PostForm("receiverId"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
        return
    }

    text := c.PostForm("text")
    if text == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Message text is required"})
        return
    }

    // Create the message
    message := models.Message{
        SenderID:   user.ID,
        ReceiverID: uint(receiverID),
        Text:       text,
    }

    // Save the message to the database
    if err := database.DB.Create(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
        return
    }

    // Create a notification for the receiver
    notification := models.Notification{
        UserID:    uint(receiverID),
        MessageID: message.ID,
        SenderID:  user.ID,
        IsRead:    false,
    }

    // Save the notification to the database
    if err := database.DB.Create(&notification).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
        return
    }

    // Render the new message using a template
    templates.SingleMessage(message, true).Render(context.Background(), c.Writer)
}

