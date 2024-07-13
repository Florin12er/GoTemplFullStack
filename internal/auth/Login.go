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
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    if err := c.BindJSON(&body); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "failed to read body",
        })
        return
    }

    var user models.User
    if err := database.DB.First(&user, "email = ?", body.Email).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid email or password",
        })
        return
    }

    if user.ID == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid email or password",
        })
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid email or password",
        })
        return
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "sub": user.ID,
        "exp": time.Now().Add(time.Hour).Unix(),
    })
    tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "failed to create token",
        })
        return
    }

    c.SetSameSite(http.SameSiteLaxMode)
    c.SetCookie("Authorization", tokenString, 3600, "", "", false, true)


    c.JSON(http.StatusOK, gin.H{})
}
