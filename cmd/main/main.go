package main

import (
	"GoMessageApp/internal/controllers"
	"GoMessageApp/internal/templates"
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Serving static files
	r.Static("/static", "./static")

	// Define a route for the main page or component
	r.GET("/", func(c *gin.Context) {
		// Create the component
		component := templates.Hello("John", "Berlin", "title")

		// Render the component to the response writer and handle any potential errors
		if err := component.Render(context.Background(), c.Writer); err != nil {
			log.Printf("failed to render component: %v", err)
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
    r.GET("/hello",controllers.HelloHandler)
	})

	// Start the Gin server and log any errors if it fails to start
    if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

