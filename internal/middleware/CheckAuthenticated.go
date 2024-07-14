package middleware

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"fmt"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)


func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key to validate the token
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !token.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Check if token is expired
	if time.Now().Unix() > int64(expFloat) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Fetch user from database using sub (subject) claim
	sub, ok := claims["sub"].(float64)  // Assuming sub is an ID and stored as a float64
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", int64(sub)).Error; err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if user.ID == 0 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Set user context for subsequent handlers
	c.Set("user", user)

	c.Next()
}

