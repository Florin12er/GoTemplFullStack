package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode"
)

// RegisterRequest is the struct used for binding and validation
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func isPasswordValid(password string) bool {
	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return len(password) >= 6 && hasUpper && hasLower && hasNumber && hasSpecial
}

func Register(c *gin.Context) {
	var body RegisterRequest

	// Bind the JSON body
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body or missing required fields",
		})
		return
	}

	// Validate password
	if !isPasswordValid(body.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password doesn't meet the requirements. It must be at least 6 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.",
		})
		return
	}

	// Generate a hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
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
			"error": "Failed to create user",
		})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

