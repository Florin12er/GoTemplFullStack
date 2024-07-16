// auth/logout.go

package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	// Clear the cookie
	c.SetCookie(
		"Authorization",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
    c.Redirect(http.StatusSeeOther, "/login")
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

