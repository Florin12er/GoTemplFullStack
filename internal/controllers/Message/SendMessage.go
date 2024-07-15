package message

import(
"github.com/gin-gonic/gin"
    "net/http"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/Database"
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


