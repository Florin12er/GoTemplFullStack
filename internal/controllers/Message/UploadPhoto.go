package message

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
)

func UploadMessagePicture(c *gin.Context) {
    // Get the user ID from the JWT token
    userID, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
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
    if message.SenderID != uint(userID.(uint)) {
        c.JSON(http.StatusForbidden, gin.H{"error": "You can only add pictures to your own messages"})
        return
    }

    // Get the file from the request
    file, err := c.FormFile("picture")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Validate file type (optional, but recommended)
    if !isValidImageType(file.Header.Get("Content-Type")) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed."})
        return
    }

    // Generate a unique filename
    filename := fmt.Sprintf("message_%d_%s", messageID, file.Filename)
    filepath := filepath.Join("uploads", "messages", filename)

    // Save the file
    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Update the message with the picture path
    message.Picture = filepath
    if err := database.DB.Save(&message).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Picture uploaded successfully", "filepath": filepath})
}

// Helper function to validate image types
func isValidImageType(contentType string) bool {
    validTypes := []string{"image/jpeg", "image/png", "image/gif"}
    for _, t := range validTypes {
        if t == contentType {
            return true
        }
    }
    return false
}

