package auth

import(
    "net/http"
    "github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
    user, _ := c.Get("user")
    c.JSON(http.StatusOK, gin.H{
        "message": user,
    })
}
