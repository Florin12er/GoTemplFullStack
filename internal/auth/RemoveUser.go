package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeleteUserProfile(c *gin.Context) {
	userID := utils.GetLoggedInUserID(c)

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
