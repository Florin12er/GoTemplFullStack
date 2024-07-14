package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadUserProfilePicture handles uploading a profile picture for a user
func UploadProfilePicture(c *gin.Context) {
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

    // Validate file type (optional, but recommended)
    if !isValidImageType(file.Header.Get("Content-Type")) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed."})
        return
    }

    // Generate a unique filename
    filename := fmt.Sprintf("%d_%s", userID, file.Filename)
    filepath := filepath.Join("uploads/users", filename)

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

    // If there's an existing profile picture, delete it
    if user.ProfilePicture != "" {
        if err := os.Remove(user.ProfilePicture); err != nil {
            // Log the error, but don't stop the process
            fmt.Printf("Failed to delete old profile picture: %v\n", err)
        }
    }

    user.ProfilePicture = filepath
    if err := database.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded successfully", "filepath": filepath})
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

