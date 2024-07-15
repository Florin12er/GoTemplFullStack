package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
    "GoMessageApp/internal/templates"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"unicode"
    "context"
)

// check if the password is valid
func isPasswordValid(password string) bool {
	if len(password) < 6 {
		return false
	}

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

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// RegisterRequest is the struct used for binding and validation
type RegisterRequest struct {
	Username string `form:"username" binding:"required"`
	Fullname string `form:"fullname" binding:"required"`
	Email    string `form:"email"    binding:"required,email"`
	Password string `form:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
    var body RegisterRequest

    // Bind form data instead of JSON
    if err := c.Bind(&body); err != nil {
        templates.RegisterForm("Failed to read form data or missing required fields").Render(context.Background(), c.Writer)
        return
    }

    // Validate password
    if !isPasswordValid(body.Password) {
        templates.RegisterForm("Password must be at least 6 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character").Render(context.Background(), c.Writer)
        return
    }

    // Generate a hash of the password
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
    if err != nil {
        templates.RegisterForm("Failed to process password").Render(context.Background(), c.Writer)
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
        templates.RegisterForm("Failed to create user. Email or username may already be in use.").Render(context.Background(), c.Writer)
        return
    }

    // Respond with a success message
    c.Header("HX-Redirect", "/login")
    c.String(http.StatusOK, "")
}
