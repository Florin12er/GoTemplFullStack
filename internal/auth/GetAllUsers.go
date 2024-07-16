package auth

import (
    "GoMessageApp/internal/Database"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/templates"
    "context"
    "github.com/gin-gonic/gin"
)

func GetAllUsers(c *gin.Context) {
    var users []models.User

    // Query the database for all users
    if err := database.DB.Find(&users).Error; err != nil {
        templates.Error("Failed to retrieve users").Render(context.Background(), c.Writer)
        return
    }

    // Render the users list template
    templates.UsersList(users).Render(context.Background(), c.Writer)
}

