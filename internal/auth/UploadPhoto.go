package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
    "log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)
func UploadProfilePicture(c *gin.Context) {
    // Get the user from the context
    userInterface, exists := c.Get("user")
    if !exists {
        c.HTML(http.StatusUnauthorized, "error.templ", gin.H{"error": "User not authenticated"})
        return
    }
    currentUser, ok := userInterface.(models.User)
    if !ok {
        c.HTML(http.StatusInternalServerError, "error.templ", gin.H{"error": "Invalid user data"})
        return
    }

    // Get the file from the request
    file, header, err := c.Request.FormFile("profilePicture")
    if err != nil {
        log.Printf("Error getting file: %v", err)
        c.HTML(http.StatusBadRequest, "error.templ", gin.H{"error": fmt.Sprintf("No file uploaded: %v", err)})
        return
    }
    defer file.Close()

    // Generate a unique filename
    ext := filepath.Ext(header.Filename)
    filename := fmt.Sprintf("%d_%s%s", currentUser.ID, uuid.New().String(), ext)
    
    // Ensure the upload directory exists
    uploadDir := "uploads/users"
    if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
        log.Printf("Error creating directory: %v", err)
        c.HTML(http.StatusInternalServerError, "error.templ", gin.H{"error": fmt.Sprintf("Failed to create upload directory: %v", err)})
        return
    }

    // Save the file
    filepath := filepath.Join(uploadDir, filename)
    if err := c.SaveUploadedFile(header, filepath); err != nil {
        log.Printf("Error saving file: %v", err)
        c.HTML(http.StatusInternalServerError, "error.templ", gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
        return
    }

    // Update user's profile picture in the database
    relativePath := path.Join("/uploads/users", filename) // Use forward slashes for URL
    if err := database.DB.Model(&currentUser).Update("profile_picture", relativePath).Error; err != nil {
        log.Printf("Error updating user in database: %v", err)
        c.HTML(http.StatusInternalServerError, "error.templ", gin.H{"error": fmt.Sprintf("Failed to update user profile: %v", err)})
        return
    }

    // Return success
    c.HTML(http.StatusOK, "profile_picture.templ", gin.H{
        "message": "Profile picture updated successfully",
        "path": relativePath,
    })
}


// EditProfilePicture updates the profile picture of the logged-in user
func EditProfilePicture(c *gin.Context) {
    // Get the user ID from the JWT token
    userID, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Get the file from the request
    file, err := c.FormFile("profilePicture")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Generate a unique filename
    filename := fmt.Sprintf("%d_%s", userID, file.Filename)
    filepath := filepath.Join("uploads/users/", filename)

    // Save the file
    if err := c.SaveUploadedFile(file, filepath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Update the user's profile picture in the database
    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    user.ProfilePicture = filepath
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile picture updated successfully", "filepath": filepath})
}

// DeleteProfilePicture deletes the profile picture of the logged-in user
func DeleteProfilePicture(c *gin.Context) {
    // Get the user ID from the JWT token
    userID, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    // Fetch the user from the database
    var user models.User
    if err := database.DB.First(&user, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Check if the user has a profile picture
    if user.ProfilePicture == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No profile picture to delete"})
        return
    }

    // Delete the file from the server
    if err := os.Remove(user.ProfilePicture); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
        return
    }

    // Update the user's profile in the database
    user.ProfilePicture = ""
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile picture deleted successfully"})
}
