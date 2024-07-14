package auth

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/models"
	"crypto/rand"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResetCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

var resetCodes = make(map[string]ResetCode)
func generateResetCode() string {
    code := make([]byte, 3)
    rand.Read(code)
    return fmt.Sprintf("%06d", int(code[0])<<16|int(code[1])<<8|int(code[2]))
}
func sendResetEmail(email, code string) error {
    from := os.Getenv("EMAIL_ADDRESS")
    password := os.Getenv("EMAIL_PASSWORD")
    to := []string{email}
    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    // HTML email template
    htmlTemplate := `
    <!DOCTYPE html>
    <html lang="en">
    <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1.0">
        <title>Password Reset</title>
        <style>
            body {
                font-family: Arial, sans-serif;
                line-height: 1.6;
                color: #333;
                max-width: 600px;
                margin: 0 auto;
                padding: 20px;
            }
            .container {
                background-color: #f9f9f9;
                border-radius: 5px;
                padding: 20px;
                text-align: center;
            }
            h1 {
                color: #2c3e50;
            }
            .code {
                font-size: 36px;
                font-weight: bold;
                color: #3498db;
                letter-spacing: 5px;
                margin: 20px 0;
                padding: 10px;
                background-color: #ecf0f1;
                border-radius: 5px;
            }
            .footer {
                margin-top: 20px;
                font-size: 12px;
                color: #7f8c8d;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <h1>Password Reset</h1>
            <p>You have requested to reset your password. Use the following code to complete the process:</p>
            <div class="code">%s</div>
            <p>This code will expire in 15 minutes.</p>
            <p>If you did not request a password reset, please ignore this email or contact support if you have concerns.</p>
            <div class="footer">
                <p>This is an automated message, please do not reply to this email.</p>
            </div>
        </div>
    </body>
    </html>
    `

    // Format the HTML template with the reset code
    htmlBody := fmt.Sprintf(htmlTemplate, code)

    // Compose the email
    mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
    subject := "Subject: Password Reset\n"
    msg := []byte(subject + mime + htmlBody)

    // Authenticate and send the email
    auth := smtp.PlainAuth("", from, password, smtpHost)
    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
    if err != nil {
        return err
    }

    return nil
}

func ResetRequest(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the email exists in the database
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	// Generate and store reset code
	code := generateResetCode()
	resetCodes[req.Email] = ResetCode{
		Email:     req.Email,
		Code:      code,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	// Send email with reset code
	if err := sendResetEmail(req.Email, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send reset email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset code sent to email"})
}

func ResetPassword(c *gin.Context) {
    var req struct {
        Email       string `json:"email" binding:"required,email"`
        Code        string `json:"code" binding:"required,min=6"`
        NewPassword string `json:"newPassword" binding:"required,min=6"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Verify reset code
    resetCode, exists := resetCodes[req.Email]
    if !exists {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No reset code found for this email"})
        return
    }
    if resetCode.Code != req.Code {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reset code"})
        return
    }
    if time.Now().After(resetCode.ExpiresAt) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Reset code has expired"})
        return
    }

    // Update password in the database
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }

    if err := database.DB.Model(&models.User{}).Where("email = ?", req.Email).Update("password", string(hashedPassword)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    // Remove used reset code
    delete(resetCodes, req.Email)

    c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}

