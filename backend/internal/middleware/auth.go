package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AuthMiddleware(sessionManager *scs.SessionManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !sessionManager.Exists(c.Request.Context(), "userID") {
			zap.S().Warnw("Unauthorized access attempt", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
