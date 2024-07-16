package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/templates"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadProfilePicture(c *gin.Context) {
    // Log request details
    fmt.Printf("Content-Type: %s\n", c.ContentType())
    fmt.Printf("Request Method: %s\n", c.Request.Method)

    // Get the user from the context
    userInterface, exists := c.Get("user")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    currentUser, ok := userInterface.(models.User)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user data"})
        return
    }

    // Get the file from the request
    file, err := c.FormFile("profilePicture")
    if err != nil {
        fmt.Printf("Error getting file: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("No file uploaded: %v", err)})
        return
    }

    // Validate file type (optional, but recommended)
    if !isValidImageType(file.Header.Get("Content-Type")) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only images are allowed."})
        return
    }

    // Generate a unique filename
    filename := fmt.Sprintf("%d_%s", currentUser.ID, file.Filename)
    uploadPath := filepath.Join("uploads", "users")
    filePath := filepath.Join(uploadPath, filename)

    // Ensure the upload directory exists
    if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
        fmt.Printf("Error creating directory: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create upload directory: %v", err)})
        return
    }

    // Save the file
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        fmt.Printf("Error saving file: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save file: %v", err)})
        return
    }

    // If there's an existing profile picture, delete it
    if currentUser.ProfilePicture != "" {
        oldFilePath := filepath.Join(".", currentUser.ProfilePicture)
        if err := os.Remove(oldFilePath); err != nil {
            fmt.Printf("Failed to delete old profile picture: %v\n", err)
            // Log the error, but don't stop the process
        }
    }

    // Update the user's profile picture in the database
    currentUser.ProfilePicture = "/" + filepath.Join("uploads", "users", filename) // Use forward slash for URL path
    if err := database.DB.Save(&currentUser).Error; err != nil {
        fmt.Printf("Error updating user in database: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update user: %v", err)})
        return
    }

    // Log success
    fmt.Printf("Profile picture uploaded successfully for user %d\n", currentUser.ID)

    // Render the updated profile content
    if err := templates.UserProfileContent(currentUser).Render(context.Background(), c.Writer); err != nil {
        fmt.Printf("Error rendering template: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to render profile: %v", err)})
        return
    }
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

