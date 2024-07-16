package content 

import (
    "GoMessageApp/internal/Database"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/templates"
    "github.com/gin-gonic/gin"
    "context"
)

func DashboardHandler(c *gin.Context) {
    userInterface, _ := c.Get("user")
    currentUser, _ := userInterface.(models.User)

    var users []models.User
    if err := database.DB.Where("id != ?", currentUser.ID).Find(&users).Error; err != nil {
        // Handle error
        templates.Error("Failed to fetch users").Render(context.Background(), c.Writer)
        return
    }

    templates.DashBoard(currentUser, users).Render(context.Background(), c.Writer)
}

