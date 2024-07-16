// internal/controllers/content/profile.go
package auth

import (
    "GoMessageApp/internal/Database"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/templates"
    "context"
    "net/http"

    "github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
    userInterface, _ := c.Get("user")
    currentUser, _ := userInterface.(models.User)

    // Fetch the latest user data from the database
    if err := database.DB.First(&currentUser, currentUser.ID).Error; err != nil {
        c.String(http.StatusInternalServerError, "Error fetching user data")
        return
    }

    templates.UserProfileContent(currentUser).Render(context.Background(), c.Writer)
}

