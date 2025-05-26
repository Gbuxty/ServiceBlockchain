package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (h *QuotesHandlers) BasicAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok || username != h.cfg.User || password != h.cfg.Password {
			c.AbortWithStatusJSON(http.StatusUnauthorized, QuitResponse{
				Success: false,
				Message: "Unauthorized access",
			})
			return
		}
		c.Next()
	}
}
