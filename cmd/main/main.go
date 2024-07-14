package main

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/auth"
	"GoMessageApp/internal/controllers"
	"GoMessageApp/internal/controllers/Message"
	"GoMessageApp/internal/controllers/notification"
	"GoMessageApp/internal/middleware"
	"GoMessageApp/internal/templates"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// load the env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
	r.GET("/", func(c *gin.Context) {
		// Create the component
		component := templates.Hello("John", "Berlin", "title")
		// Render the component to the response writer and handle any potential errors
		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failed to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	})
	// Authenticate user
	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	r.POST("/reset", auth.ResetPassword)
	r.GET("/user", middleware.RequireAuth, auth.GetUser)
	r.GET("/user/all", middleware.RequireAuth, auth.GetAllUsers)
	r.GET("/user/:userID", middleware.RequireAuth, auth.GetUserByID)
	r.PUT("/user/profile", middleware.RequireAuth, auth.EditUserProfile)
	r.DELETE("/user/profile", middleware.RequireAuth, auth.DeleteUserProfile)

	// Messages
	r.POST("/message", middleware.RequireAuth, message.SendMessage)
	r.PUT("/message/:messageID", middleware.RequireAuth, message.EditMessage)
	r.DELETE("/message/:messageID", middleware.RequireAuth, message.DeleteMessage)
	r.GET("/conversations/user/:userID", middleware.RequireAuth, message.GetAllConversations)
	r.GET("/conversation/:senderID/:receiverID", middleware.RequireAuth, message.GetConversation)

	// Notification

	r.GET("/notification", middleware.RequireAuth, notification.GetNotifications)
	r.DELETE(
		"/notification/:notificationID",
		middleware.RequireAuth,
		notification.RemoveNotification,
	)

	// Upload photos
	r.POST("/user/:userID/profile-picture", middleware.RequireAuth, auth.UploadUserProfilePicture)
	r.POST("/user/profile-picture", middleware.RequireAuth, auth.EditProfilePicture)
	r.DELETE("/user/profile-picture", middleware.RequireAuth, auth.DeleteProfilePicture)
	r.POST("/messages/:messageID/picture", middleware.RequireAuth, message.UploadMessagePicture)

	// function routes
	r.GET("/hello", controllers.HelloHandler)

	// Run the server on the specified port
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
