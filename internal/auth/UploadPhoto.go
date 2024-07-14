package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadUserProfilePicture handles uploading a profile picture for a user
func UploadUserProfilePicture(c *gin.Context) {
	userID := c.Param("userID")

	// Retrieve the user from the database
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Handle file upload
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no file uploaded"})
		return
	}

	// Generate a unique filename for the uploaded file
	filename := filepath.Base(file.Filename)
	filepath := "uploads/users/" + filename

	// Save the uploaded file to the specified filepath
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Update the user's profile picture URL
	user.ProfilePicture = filepath
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{"message": "profile picture uploaded successfully", "url": filepath},
	)
}

// EditProfilePicture updates the profile picture of the logged-in user
func EditProfilePicture(c *gin.Context) {
	userID := utils.GetLoggedInUserID(
		c,
	) // Assume getLoggedInUserID is implemented as discussed earlier

	file, err := c.FormFile("profile_picture")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing profile picture file"})
		return
	}

	// Generate a unique filename for the uploaded file
	filename := filepath.Base(file.Filename)
	filepath := fmt.Sprintf("uploads/users/%d_%s", userID, filename)

	// Save the file to the specified filepath
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save profile picture"})
		return
	}

	// Update the user record with the new profile picture path
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	user.ProfilePicture = filepath // Update with the actual filepath where the picture is stored

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile picture updated successfully"})
}

// DeleteProfilePicture deletes the profile picture of the logged-in user
func DeleteProfilePicture(c *gin.Context) {
	userID := utils.GetLoggedInUserID(
		c,
	) // Assume getLoggedInUserID is implemented as discussed earlier

	// Fetch the existing user record
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Check if user already has a profile picture
	if user.ProfilePicture == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile picture does not exist"})
		return
	}

	// Delete the profile picture file from disk
	if err := os.Remove(user.ProfilePicture); err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to delete profile picture file"},
		)
		return
	}

	// Clear the profile picture path in the user record
	user.ProfilePicture = ""

	// Save the updated user record
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile picture deleted successfully"})
}
