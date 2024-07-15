package message

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

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


