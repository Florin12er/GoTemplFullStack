package main

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/auth"
	"GoMessageApp/internal/controllers"
	"GoMessageApp/internal/controllers/Message"
	"GoMessageApp/internal/controllers/notification"
	"GoMessageApp/internal/middleware"
	"GoMessageApp/internal/templates"
	"GoMessageApp/internal/utils"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	r.GET("/", middleware.RequireAuth,func(c *gin.Context) {
		// Create the component
		component := templates.DashBoard("John", "Berlin")
		// Render the component to the response writer and handle any potential errors
		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failed to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	})
	r.GET("/login",middleware.RedirectIfAuthenticated(), func(c *gin.Context) {
		component := templates.Login()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})
	r.GET("/register", middleware.RedirectIfAuthenticated(),func(c *gin.Context) {
		component := templates.Register()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})
	r.GET("/reset-request",middleware.RedirectIfAuthenticated() ,func(c *gin.Context) {
		component := templates.ResetRequest()

		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failded to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Sever Error")
		}
	})
	r.GET("/reset-password",middleware.RedirectIfAuthenticated(), func(c *gin.Context) {
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
	r.GET("/auth/user", middleware.RequireAuth, auth.GetUser)
	r.GET("/auth/user/all", middleware.RequireAuth, auth.GetAllUsers)
	r.GET("/auth/user/:userID", middleware.RequireAuth, auth.GetUserByID)
	r.PUT("/auth/user/profile", middleware.RequireAuth, auth.EditUserProfile)
	r.DELETE("/auth/user/profile", middleware.RequireAuth, auth.DeleteUserProfile)

	// Messages
	r.POST("/auth/message", middleware.RequireAuth, message.SendMessage)
	r.PUT("/auth/message/:messageID", middleware.RequireAuth, message.EditMessage)
	r.DELETE("/auth/message/:messageID", middleware.RequireAuth, message.DeleteMessage)
	r.GET("/auth/conversation/:receiverID", middleware.RequireAuth, message.GetConversation)
	r.GET("/auth/conversations", middleware.RequireAuth, message.GetAllConversations)

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

	// function routes
	r.GET("/hello", controllers.HelloHandler)

	// Run the server on the specified port
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
