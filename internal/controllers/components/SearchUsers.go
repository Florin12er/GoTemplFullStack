// handlers/search_users.go
package components

import (
    "GoMessageApp/internal/Database"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/templates"
    "github.com/gin-gonic/gin"
    "strings"
)

func SearchUsers(c *gin.Context) {
    query := strings.ToLower(c.PostForm("user"))

    var users []models.User
    if query == "" {
        // If query is empty, return all users
        if err := database.DB.Find(&users).Error; err != nil {
            templates.Error("Failed to fetch users").Render(c, c.Writer)
            return
        }
    } else {
        // If there's a query, filter users
        if err := database.DB.Where("LOWER(user_name) LIKE ? OR LOWER(full_name) LIKE ?", "%"+query+"%", "%"+query+"%").Find(&users).Error; err != nil {
            templates.Error("Failed to search users").Render(c, c.Writer)
            return
        }
    }

    templates.UsersList(users).Render(c, c.Writer)
}

