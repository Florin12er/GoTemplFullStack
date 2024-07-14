package message

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// SendMessage allows a user to send a message to another user
func SendMessage(c *gin.Context) {
    var body struct {
        SenderID   uint   `json:"sender_id"`
        ReceiverID uint   `json:"receiver_id"`
        Text       string `json:"text"`
    }

    if err := c.BindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    message := models.Message{
        SenderID:   body.SenderID,
        ReceiverID: body.ReceiverID,
        Text:       body.Text,
    }

    if err := database.DB.Create(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not send message"})
        return
    }

    // Create notification for the receiver
    notification := models.Notification{
        UserID:    body.ReceiverID,
        MessageID: message.ID,
        SenderID:  body.SenderID,
        IsRead:    false, // By default, notification is unread
    }
    if err := database.DB.Create(&notification).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create notification"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "message sent successfully"})
}

// GetConversation retrieves all messages between two users
func GetConversation(c *gin.Context) {
    senderID, err := strconv.Atoi(c.Param("senderID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sender ID"})
        return
    }

    receiverID, err := strconv.Atoi(c.Param("receiverID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiver ID"})
        return
    }

    var messages []models.Message
    if err := database.DB.Preload("Sender").Preload("Receiver").
        Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
            senderID, receiverID, receiverID, senderID).
        Order("created_at").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve conversation"})
        return
    }

    c.JSON(http.StatusOK, messages)
}

// GetAllConversations retrieves all conversations for a user
func GetAllConversations(c *gin.Context) {
    userID, err := strconv.Atoi(c.Param("userID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    var messages []models.Message
    if err := database.DB.Preload("Sender").Preload("Receiver").
        Where("sender_id = ? OR receiver_id = ?", userID, userID).
        Order("created_at").Find(&messages).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not retrieve conversations"})
        return
    }

    c.JSON(http.StatusOK, messages)
}
// EditMessage allows a user to edit a message
func EditMessage(c *gin.Context) {
    messageID, err := strconv.Atoi(c.Param("messageID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
        return
    }

    var update struct {
        Text string `json:"text"`
    }
    if err := c.BindJSON(&update); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Fetch the message
    var message models.Message
    if err := database.DB.First(&message, messageID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
        return
    }

    // Update the message text
    message.Text = update.Text
    if err := database.DB.Save(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update message"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "message updated successfully"})
}

// DeleteMessage allows a user to delete a message
func DeleteMessage(c *gin.Context) {
    messageID, err := strconv.Atoi(c.Param("messageID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
        return
    }

    // Fetch the message
    var message models.Message
    if err := database.DB.First(&message, messageID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
        return
    }

    // Soft delete the message (assuming you have a DeletedAt field in your Message model)
    if err := database.DB.Delete(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete message"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "message deleted successfully"})
}

