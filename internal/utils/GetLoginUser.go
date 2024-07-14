package utils

import (
	"github.com/gin-gonic/gin"
)
// Assuming this function is defined in your middleware or auth package
func GetLoggedInUserID(c *gin.Context) uint {
    userID, exists := c.Get("userID")
    if !exists {
        return 0 // Or handle as needed, like returning an error
    }
    if id, ok := userID.(uint); ok {
        return id
    }
    return 0 // Or handle as needed
}

