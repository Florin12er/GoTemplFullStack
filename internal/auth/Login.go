package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"GoMessageApp/internal/templates"
	"context"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
    var body struct {
        Email    string `form:"email"`
        Password string `form:"password"`
    }
    if err := c.Bind(&body); err != nil {
        templates.Error("Failed to read body").Render(context.Background(), c.Writer)
        return
    }

    var user models.User
    if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
        templates.Error("Invalid email or password").Render(context.Background(), c.Writer)
        return
    }

    if user.ID == 0 {
        templates.Error("Invalid email or password").Render(context.Background(), c.Writer)
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
        templates.Error("Invalid email or password").Render(context.Background(), c.Writer)
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour).Unix(),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
    if err != nil {
        templates.Error("Failed to create token").Render(context.Background(), c.Writer)
        return
    }

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

    c.Header("HX-Redirect", "/")
    c.String(http.StatusOK, "")
}

