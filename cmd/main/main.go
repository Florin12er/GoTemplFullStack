package main

import (
	"GoMessageApp/cmd/main/websocket"
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/auth"
	"GoMessageApp/internal/controllers/Message"
	"GoMessageApp/internal/controllers/components"
	"GoMessageApp/internal/controllers/content"
	"GoMessageApp/internal/controllers/notification"
	"GoMessageApp/internal/middleware"
	"GoMessageApp/internal/templates"
	"GoMessageApp/internal/utils"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	utils.LoadEnv()
	database.ConnectToDb()
	database.SyncDatabase()
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	// initialize gin
	r := gin.Default()

	if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
		panic(err)
	}

	// Serving static files
	r.Static("/static", "./static")
	r.Static("/uploads", "./uploads")

	//html rendering routes
	r.GET("/", middleware.RequireAuth, content.DashboardHandler)
	r.GET("/login", middleware.RedirectIfAuthenticated(), func(c *gin.Context) {
		component := templates.Login()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})
	r.GET("/register", middleware.RedirectIfAuthenticated(), func(c *gin.Context) {
		component := templates.Register()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})
	r.GET("/reset-request", middleware.RedirectIfAuthenticated(), func(c *gin.Context) {
		component := templates.ResetRequest()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})
	r.GET("/reset-password", middleware.RedirectIfAuthenticated(), func(c *gin.Context) {
		component := templates.ResetPassword()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})

	// Authenticate user
	r.POST("/auth/register", auth.Register)
	r.POST("/auth/login", auth.Login)
	r.POST("/auth/logout", middleware.RequireAuth, auth.Logout)
	r.POST("/auth/reset-request", auth.ResetRequest)
	r.POST("/auth/reset-password", auth.ResetPassword)
	r.GET("/profile-content", middleware.RequireAuth, auth.GetUser)
	r.GET("/auth/user/all", middleware.RequireAuth, auth.GetAllUsers)
	r.POST("/search-users", middleware.RequireAuth, components.SearchUsers)
	r.GET("/auth/user/:userID", middleware.RequireAuth, auth.GetUserByID)
	r.PUT("/auth/user/profile", middleware.RequireAuth, auth.EditProfile)
	r.DELETE("/auth/user/profile", middleware.RequireAuth, auth.DeleteUserProfile)
	r.PUT("/auth/total/user/profile", middleware.RequireAuth, auth.EditUserProfile)

	// Messages
	r.POST("/auth/message", middleware.RequireAuth, message.SendMessage)
	r.GET("/ws", middleware.RequireAuth, websocket.HandleWebSocket)
	r.PUT("/auth/message/:messageID", middleware.RequireAuth, message.EditMessage)
	r.DELETE("/auth/message/:messageID", middleware.RequireAuth, message.DeleteMessage)
	r.GET("/conversation/:userID", middleware.RequireAuth, message.GetConversation)
	r.GET("/conversations", middleware.RequireAuth, message.GetAllConversations)

	// Notification

	r.GET("/auth/notification", middleware.RequireAuth, notification.GetNotifications)
	r.DELETE(
		"/auth/notification/:notificationID",
		middleware.RequireAuth,
		notification.RemoveNotification,
	)

	// Upload photos
	r.POST("/auth/user/:userID/profile-picture", middleware.RequireAuth, auth.UploadProfilePicture)
	r.POST("/auth/user/profile-picture", middleware.RequireAuth, auth.EditProfilePicture)
	r.DELETE("/auth/user/profile-picture", middleware.RequireAuth, auth.DeleteProfilePicture)
	r.POST(
		"/auth/messages/:messageID/picture",
		middleware.RequireAuth,
		message.UploadMessagePicture,
	)

	// Run the server on the specified port
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
