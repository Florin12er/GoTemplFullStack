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
func GetConversation(c *gin.Context) {
    userInterface, _ := c.Get("user")
    currentUser, _ := userInterface.(models.User)

    receiverID, err := strconv.ParseUint(c.Param("userID"), 10, 32)
    if err != nil {
        templates.Error("Invalid user ID").Render(context.Background(), c.Writer)
        return
    }

    var receiver models.User
    if err := database.DB.First(&receiver, receiverID).Error; err != nil {
        templates.Error("User not found").Render(context.Background(), c.Writer)
        return
    }

    var messages []models.Message
    if err := database.DB.Where(
        "(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
        currentUser.ID, receiverID, receiverID, currentUser.ID,
    ).Order("created_at ASC").Find(&messages).Error; err != nil {
        templates.Error("Failed to fetch messages").Render(context.Background(), c.Writer)
        return
    }

    templates.Conversation(currentUser, receiver, messages).Render(context.Background(), c.Writer)
}

func GetAllConversations(c *gin.Context) {
	userInterface, _ := c.Get("user")
	currentUser, _ := userInterface.(models.User)

	var conversations []templates.ConversationPreview

	// Query to get all unique conversations with last message
	query := `
        SELECT 
            CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS user_id,
            CASE WHEN sender_id = ? THEN receiver.user_name ELSE sender.user_name END AS user_name,
            m.text AS last_message,
            m.created_at AS last_message_time
        FROM messages m
        JOIN users AS sender ON sender.id = m.sender_id
        JOIN users AS receiver ON receiver.id = m.receiver_id
        INNER JOIN (
            SELECT 
                CASE WHEN sender_id = ? THEN receiver_id ELSE sender_id END AS user_id,
                MAX(created_at) AS max_created_at
            FROM messages
            WHERE sender_id = ? OR receiver_id = ?
            GROUP BY user_id
        ) latest ON (latest.user_id = CASE WHEN m.sender_id = ? THEN m.receiver_id ELSE m.sender_id END)
            AND m.created_at = latest.max_created_at
        WHERE m.sender_id = ? OR m.receiver_id = ?
        ORDER BY m.created_at DESC
    `

	if err := database.DB.Raw(query, currentUser.ID, currentUser.ID, currentUser.ID, currentUser.ID, currentUser.ID, currentUser.ID, currentUser.ID, currentUser.ID).Scan(&conversations).Error; err != nil {
		c.HTML(
			http.StatusInternalServerError,
			"error.html",
			gin.H{"error": "Failed to fetch conversations"},
		)
		return
	}

	templates.ConversationsList(conversations).Render(context.Background(), c.Writer)
}
