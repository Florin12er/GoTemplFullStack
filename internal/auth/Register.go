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

    // Bind the form data
    if err := c.ShouldBind(&body); err != nil {
        c.HTML(http.StatusBadRequest, "", "Failed to read form data or missing required fields")
        return
    }

    // Check if username is already taken
    var existingUser models.User
    if err := database.DB.Where("user_name = ?", body.Username).First(&existingUser).Error; err == nil {
        c.HTML(http.StatusBadRequest, "", "Username is already taken")
        return
    }

    // Validate password
    if !isPasswordValid(body.Password) {
        c.HTML(http.StatusBadRequest, "", "Password doesn't meet the requirements. It must be at least 6 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.")
        return
    }

    // Generate a hash of the password
    hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
    if err != nil {
        c.HTML(http.StatusInternalServerError, "", "Failed to hash password")
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
        c.HTML(http.StatusInternalServerError, "", "Failed to create user")
        return
    }

    // Respond with a success message
    c.HTML(http.StatusOK, "", "<p>User created successfully! <a href='/login'>Login here</a></p>")
}
func CheckUsername(c *gin.Context) {
    username := c.PostForm("username")
    
    var existingUser models.User
    if err := database.DB.Where("user_name = ?", username).First(&existingUser).Error; err == nil {
        c.HTML(http.StatusOK, "", "Username is already taken")
    } else {
        c.HTML(http.StatusOK, "", "")
    }
}

func CheckPassword(c *gin.Context) {
    password := c.PostForm("password")
    
    if !isPasswordValid(password) {
        c.HTML(http.StatusOK, "", "Password must be at least 6 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character.")
    } else {
        c.HTML(http.StatusOK, "", "")
    }
}

