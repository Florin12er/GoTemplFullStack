package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// RegisterRequest is the struct used for binding and validation
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var body RegisterRequest

	// Bind and validate the JSON body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body or missing required fields",
		})
		return
	}

	// Generate a hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to hash password",
		})
		return
	}

	// Create a new user instance
	user := models.User{
		UserName: body.Username,
		FullName: body.Fullname,
		Email:    body.Email,
		Password: string(hash),
	}

	// Insert the user record into the database
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "user created successfully",
	})
}
