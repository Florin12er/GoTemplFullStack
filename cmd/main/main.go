package main

import (
	"GoMessageApp/internal/Database"
	"GoMessageApp/internal/auth"
	"GoMessageApp/internal/controllers"
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

	// function routes
	r.GET("/hello", controllers.HelloHandler)

	// Run the server on the specified port
	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
