package message

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// SendMessage allows a user to send a message to another user

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiverId" binding:"required"`
	Text       string `json:"text" binding:"required"`
}
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

    // Parse the request body
    var req SendMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create the message
    message := models.Message{
        SenderID:   user.ID,
        ReceiverID: req.ReceiverID,
        Text:       req.Text,
    }

    // Save the message to the database
    if err := database.DB.Create(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message"})
        return
    }

    // Create a notification for the receiver
    notification := models.Notification{
        UserID:    req.ReceiverID,
        MessageID: message.ID,
        SenderID:  user.ID,
        IsRead:    false,
    }

    // Save the notification to the database
    if err := database.DB.Create(&notification).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notification"})
        return
    }

    // Prepare the response
    response := gin.H{
        "message": "Message sent successfully",
        "data": gin.H{
            "id":         message.ID,
            "senderID":   message.SenderID,
            "receiverID": message.ReceiverID,
            "text":       message.Text,
            "createdAt":  message.CreatedAt,
            "senderName": user.UserName,
        },
    }

    c.JSON(http.StatusOK, response)
}

// GetConversation retrieves all messages between two users
func GetConversation(c *gin.Context) {
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

    // Get the receiver ID from the URL parameter
    receiverID, err := strconv.ParseUint(c.Param("receiverID"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
        return
    }

    var messages []models.Message
    if err := database.DB.Where(
        "(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
        user.ID, receiverID, receiverID, user.ID,
    ).Order("created_at ASC").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversation"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"messages": messages})
}
func GetAllConversations(c *gin.Context) {
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

    var conversations []struct {
        UserID   uint   `json:"userID"`
        UserName string `json:"userName"`
    }

    // Query to get all unique conversations
    query := `
        SELECT DISTINCT
            CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS user_id,
            CASE WHEN sender_id = ? THEN receiver.user_name ELSE sender.user_name END AS user_name
        FROM messages
        JOIN users AS sender ON sender.id = messages.sender_id
        JOIN users AS receiver ON receiver.id = messages.receiver_id
        WHERE sender_id = ? OR receiver_id = ?
    `

    if err := database.DB.Raw(query, user.ID, user.ID, user.ID, user.ID).Scan(&conversations).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch conversations"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"conversations": conversations})
}

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

