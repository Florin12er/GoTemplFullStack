package message

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

// UploadMessagePicture handles uploading a picture for a message
func UploadMessagePicture(c *gin.Context) {
	messageID := c.Param("messageID")

	// Retrieve the message from the database
	var message models.Message
	if err := database.DB.First(&message, messageID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Save the uploaded file
	filename := filepath.Base(file.Filename)
	filepath := "uploads/messages/" + filename
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Update the message's picture URL
	message.Picture = filepath
	if err := database.DB.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update message"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"message": "message picture uploaded successfully", "url": filepath},
	)
}
