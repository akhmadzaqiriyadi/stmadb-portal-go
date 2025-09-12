// internal/handler/health_handler.go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// GetHealthCheck handles the health check endpoint.
// @Summary      Show the status of server
// @Description  get the status of server
// @Tags         Health Check
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func GetHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "STMADB Portal Go Backend is running!",
	})
}