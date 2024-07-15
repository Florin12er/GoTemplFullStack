package auth

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/Database"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
    "os"
    "time"
)

func Login(c *gin.Context) {
    var body struct {
        Email    string `form:"email"`
        Password string `form:"password"`
    }
    if err := c.Bind(&body); err != nil {
        c.HTML(http.StatusBadRequest, "", "Failed to read form data")
        return
    }

    var user models.User
    if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
        c.HTML(http.StatusBadRequest, "", "Invalid email or password")
        return
    }

    if user.ID == 0 {
        c.HTML(http.StatusBadRequest, "", "Invalid email or password")
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
        c.HTML(http.StatusBadRequest, "", "Invalid email or password")
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour).Unix(),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
    if err != nil {
        c.HTML(http.StatusInternalServerError, "", "Failed to create token")
        return
    }

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)

    c.HTML(http.StatusOK, "", "<p>Login successful! Redirecting...</p><script>setTimeout(function(){ window.location.href = '/'; }, 2000);</script>")
}

