package auth

import (
    "GoMessageApp/internal/Database"
    "GoMessageApp/internal/models"
    "GoMessageApp/internal/utils"
    "GoMessageApp/internal/templates"
    "time"
    "context"
    "golang.org/x/crypto/bcrypt"
    "github.com/gin-gonic/gin"
)

func ResetRequest(c *gin.Context) {
    var req struct {
        Email string `form:"email" binding:"required,email"`
    }

    if err := c.Bind(&req); err != nil {
        templates.ResetRequestForm("Invalid email address").Render(context.Background(), c.Writer)
        return
    }

    // Check if the email exists in the database
    var user models.User
    if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
        templates.ResetRequestForm("Email not found").Render(context.Background(), c.Writer)
        return
    }

    // Generate and store reset code
    code := utils.GenerateResetCode()
    utils.ResetCodes[req.Email] = utils.ResetCode{
        Email:     req.Email,
        Code:      code,
        ExpiresAt: time.Now().Add(15 * time.Minute),
    }

    // Send email with reset code
    if err := utils.SendResetEmail(req.Email, code); err != nil {
        templates.ResetRequestForm("Failed to send reset email").Render(context.Background(), c.Writer)
        return
    }

    // Render success message
    templates.ResetRequestSuccess().Render(context.Background(), c.Writer)
}
func ResetPassword(c *gin.Context) {
    var req struct {
        Email       string `form:"email" binding:"required,email"`
        Code        string `form:"code" binding:"required,min=6"`
        NewPassword string `form:"newPassword" binding:"required,min=6"`
    }

    if err := c.Bind(&req); err != nil {
        templates.ResetPasswordForm("Invalid input. Please check your email, code, and new password.").Render(context.Background(), c.Writer)
        return
    }

    // Verify reset code
    resetCode, exists := utils.ResetCodes[req.Email]
    if !exists {
        templates.ResetPasswordForm("No reset code found for this email").Render(context.Background(), c.Writer)
        return
    }
    if resetCode.Code != req.Code {
        templates.ResetPasswordForm("Invalid reset code").Render(context.Background(), c.Writer)
        return
    }
    if time.Now().After(resetCode.ExpiresAt) {
        templates.ResetPasswordForm("Reset code has expired").Render(context.Background(), c.Writer)
        return
    }

    // Update password in the database
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        templates.ResetPasswordForm("Failed to process new password").Render(context.Background(), c.Writer)
        return
    }

    if err := database.DB.Model(&models.User{}).Where("email = ?", req.Email).Update("password", string(hashedPassword)).Error; err != nil {
        templates.ResetPasswordForm("Failed to update password").Render(context.Background(), c.Writer)
        return
    }

    // Remove used reset code
    delete(utils.ResetCodes, req.Email)

    // Render success message
    templates.ResetPasswordSuccess().Render(context.Background(), c.Writer)
}

