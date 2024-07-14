package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)
func HelloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello from Go!")
}