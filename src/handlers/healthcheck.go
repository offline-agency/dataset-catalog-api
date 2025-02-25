package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthcheckHandler responds with a simple OK status.
func HealthcheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
