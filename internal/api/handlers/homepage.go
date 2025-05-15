package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles health check requests
func HomePage(c *gin.Context) {
	html := `
    <!DOCTYPE html>
    <html>
        <head>
            <title>Hello World</title>
        </head>
        <body>
            <h1>Hello, Gin!</h1>
        </body>
    </html>
    `
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
